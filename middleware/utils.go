package middleware

import (
	"net/url"
	"strings"
)

func getBearerToken(h string) (string, bool) {
	ok := true
	components := strings.Fields(h)
	if len(components) == 2 {
		return strings.TrimSpace(components[1]), ok
	}
	return "", !ok
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
