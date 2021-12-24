package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/mstreet3/banking/entities"
	"github.com/mstreet3/banking/logger"
	"github.com/mstreet3/banking/utils"
)

type AuthMiddleware interface {
	TokenExists(http.Handler) http.Handler
	VerifyClaims(http.Handler) http.Handler
}

type DefaultAuthMiddleware struct {
	routeToRoleAccessMap routeAuthorization
	routeToClaimsMap     routeRequiredClaims
	accessToken          string
}
type Claims struct {
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	Accounts   []string `json:"accounts,omitempty"`
	CustomerId string   `json:"customer_id,omitempty"`
}

type roleAccess map[entities.Role]bool
type roleClaims map[entities.Role][]string
type routeAuthorization map[entities.AppRoute]roleAccess
type routeRequiredClaims map[entities.AppRoute]roleClaims

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

/* todo: add unit tests on isAuthorized */
func (amw *DefaultAuthMiddleware) isAuthorized(c Claims, route entities.AppRoute, vars map[string]string) bool {
	roleAuthorizedForRoute := false
	claimsVerified := false

	/* check that role is authorized for the route */
	allowed, present := amw.routeToRoleAccessMap[route][entities.Role(c.Role)]
	if present && allowed {
		roleAuthorizedForRoute = true
	}

	/* check if the route has any verifiable claims for roles */
	claimsForRole, checkRole := amw.routeToClaimsMap[route]

	if !checkRole {
		return roleAuthorizedForRoute
	}

	/* check if the current role claims to verify */
	claimsToVerify, checkClaims := claimsForRole[entities.Role(c.Role)]

	if !checkClaims {
		return roleAuthorizedForRoute
	}

	/* verify all claims for the role */
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
		accessToken, ok := getBearerToken(authHeader)

		if !ok {
			utils.WriteResponse(w, http.StatusUnauthorized, "invalid access token")
			return
		}

		amw.accessToken = accessToken
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(handler)
}

func (amw *DefaultAuthMiddleware) VerifyClaims(next http.Handler) http.Handler {
	var msg string
	handler := func(w http.ResponseWriter, r *http.Request) {
		/* get route data and build request to auth server */
		route := entities.AppRoute(mux.CurrentRoute(r).GetName())
		vars := mux.Vars(r)
		logger.Info(fmt.Sprintf("verify claims for route: %s", route))
		u := buildVerifyURL(amw.accessToken)

		/* get jwt verification response from auth server */
		response, err := http.Get(u)
		if err != nil {
			msg = fmt.Sprintf("Error sending request to auth server: " + err.Error())
			utils.WriteResponse(w, http.StatusUnauthorized, msg)
			return
		}

		/* return the error if jwt verification failed */
		if response.StatusCode != http.StatusOK {
			err = json.NewDecoder(response.Body).Decode(&msg)
			if err != nil {
				msg = fmt.Sprintf("Error while decoding response from auth server: %s", err.Error())
			}
			utils.WriteResponse(w, http.StatusUnauthorized, msg)
			return
		}

		/* status OK means token is authentic and non expired, parse the claims */
		var claims Claims
		err = json.NewDecoder(response.Body).Decode(&claims)

		if err != nil {
			msg = fmt.Sprintf("Error while decoding response from auth server: " + err.Error())
			utils.WriteResponse(w, http.StatusUnauthorized, msg)
			return
		}

		/* check if token claims authorize the request */
		authorized := amw.isAuthorized(claims, route, vars)
		if !authorized {
			utils.WriteResponse(w, http.StatusUnauthorized, "unauthorized request")
			return
		}

		/* request is authentic and authorized, serve next route in the chain */
		next.ServeHTTP(w, r)

	}
	return http.HandlerFunc(handler)
}
