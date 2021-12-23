package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/mstreet3/banking/entities"
	"github.com/mstreet3/banking/logger"
	"github.com/mstreet3/banking/utils"
)

type AuthMiddleware interface {
	TokenExists(http.Handler) http.Handler
	VerifyClaims(http.Handler) http.Handler
}

type roleAccess map[entities.Role]bool
type roleClaims map[entities.Role][]string
type routeAuthorization map[entities.AppRoute]roleAccess
type routeRequiredClaims map[entities.AppRoute]roleClaims

type Claims struct {
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	Accounts   []string `json:"accounts,omitempty"`
	CustomerId string   `json:"customer_id,omitempty"`
}

type DefaultAuthMiddleware struct {
	routeToRoleAccessMap routeAuthorization
	routeToClaimsMap     routeRequiredClaims
	accessToken          string
}

func newRouteAuthorizationMap() routeAuthorization {
	auth := make(map[entities.AppRoute]roleAccess)
	auth[entities.GetAllCustomers] = roleAccess{
		entities.CLIENT: false,
		entities.ADMIN:  true,
	}
	auth[entities.GetAllCustomersByStatus] = roleAccess{
		entities.CLIENT: false,
		entities.ADMIN:  true,
	}
	auth[entities.GetCustomerById] = roleAccess{
		entities.CLIENT: true,
		entities.ADMIN:  true,
	}
	auth[entities.NewAccount] = roleAccess{
		entities.CLIENT: false,
		entities.ADMIN:  true,
	}
	auth[entities.NewTransaction] = roleAccess{
		entities.CLIENT: true,
		entities.ADMIN:  true,
	}
	return auth

}

func newRouteRequiredClaimsMap() routeRequiredClaims {
	m := make(routeRequiredClaims)
	m[entities.GetCustomerById] = roleClaims{
		entities.CLIENT: []string{"customer_id"},
	}
	m[entities.NewTransaction] = roleClaims{
		entities.CLIENT: []string{"account_id"},
	}
	return m

}

func NewAuthMiddleware() []mux.MiddlewareFunc {
	amw := DefaultAuthMiddleware{
		routeToRoleAccessMap: newRouteAuthorizationMap(),
		routeToClaimsMap:     newRouteRequiredClaimsMap(),
	}
	return []mux.MiddlewareFunc{
		amw.TokenExists,
		amw.VerifyClaims,
	}
}

func getBearerToken(h string) string {
	components := strings.Fields(h)
	if len(components) == 2 {
		return strings.TrimSpace(components[1])
	}
	return ""
}

func buildVerifyURL(token string) string {
	u := url.URL{Host: "localhost:9000", Path: "/auth/verify", Scheme: "http"}
	q := u.Query()
	q.Add("token", token)
	u.RawQuery = q.Encode()
	return u.String()
}

func verifyAccountId(accts []string, id string) bool {
	found := false
	for _, acct := range accts {
		if acct == id {
			found = true
			break
		}
	}
	return found
}

func (amw *DefaultAuthMiddleware) isAuthorized(c Claims, route entities.AppRoute, vars map[string]string) bool {
	roleAuthorizedForRoute := false
	claimsVerified := false

	/* check that role is authorized for the route */
	allowed, present := amw.routeToRoleAccessMap[route][entities.Role(c.Role)]
	if present && allowed {
		roleAuthorizedForRoute = true
	}

	claimsForRole, checkRole := amw.routeToClaimsMap[route]

	if !checkRole {
		return roleAuthorizedForRoute
	}

	claimsToVerify, checkClaims := claimsForRole[entities.Role(c.Role)]

	if !checkClaims {
		return roleAuthorizedForRoute
	}

	for _, claim := range claimsToVerify {
		if claim == "customer_id" {
			claimsVerified = c.CustomerId == vars[claim]
		}
		if claim == "account_id" {
			claimsVerified = verifyAccountId(c.Accounts, vars[claim])
		}
	}

	return roleAuthorizedForRoute && claimsVerified
}

func (amw *DefaultAuthMiddleware) TokenExists(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		accessToken := getBearerToken(authHeader)
		if accessToken != "" {
			amw.accessToken = accessToken
			next.ServeHTTP(w, r)
		} else {
			utils.WriteResponse(w, http.StatusUnauthorized, "invalid access token")
		}

	}
	return http.HandlerFunc(handler)
}

func (amw *DefaultAuthMiddleware) VerifyClaims(next http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		route := entities.AppRoute(mux.CurrentRoute(r).GetName())
		vars := mux.Vars(r)
		logger.Info(fmt.Sprintf("verify claims for route: %s", route))
		u := buildVerifyURL(amw.accessToken)
		if response, err := http.Get(u); err != nil {
			utils.WriteResponse(w, http.StatusUnauthorized, fmt.Sprintf("Error sending request to auth server: "+err.Error()))
		} else {
			/* return the response body as an error verify failed */
			if response.StatusCode != http.StatusOK {
				var msg string
				err = json.NewDecoder(response.Body).Decode(&msg)
				if err != nil {
					msg = fmt.Sprintf("Error while decoding response from auth server: %s", err.Error())
				}
				utils.WriteResponse(w, http.StatusUnauthorized, msg)
				return
			}

			/* verify claims if status is ok*/
			var claims Claims
			if err = json.NewDecoder(response.Body).Decode(&claims); err != nil {
				utils.WriteResponse(w, http.StatusUnauthorized, fmt.Sprintf("Error while decoding response from auth server:"+err.Error()))
				return
			}
			if authorized := amw.isAuthorized(claims, route, vars); authorized {
				next.ServeHTTP(w, r)
				return
			}
			utils.WriteResponse(w, http.StatusUnauthorized, "unauthorized request")
			return
		}

	}
	return http.HandlerFunc(handler)
}
