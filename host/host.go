// Package host contains middleware for host allow request.
package host

import (
	"net/http"
	"strings"
)

// Allow middleware validates the requests are from authorized hosts.
// hosts must be hostnames separated by comma (,)
func Allow(hosts string) func(http.Handler) http.Handler {

	hs := strings.Split(hosts, ",")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for i := range hs {
				if r.Host == hs[i] {
					h.ServeHTTP(w, r)
					return
				}
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("unauthorized access host"))
		})
	}
}
