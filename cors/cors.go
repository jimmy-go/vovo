// Package cors contains middleware to enable CORS.
package cors

import (
	"net/http"
	"strings"
)

// Handler handler middleware with allow-access-control-origins header.
// origins is a string with hosts separated by comma.
func Handler(origins string) func(http.Handler) http.Handler {
	ors := strings.Split(origins, ",")
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			for _, k := range ors {
				s := struct{ domain string }{k} // copy for capturing by the goroutine
				if s.domain == r.Header.Get("Origin") || s.domain == "*" {
					w.Header().Set("Access-Control-Allow-Origin", s.domain)
					h.ServeHTTP(w, r)
					return
				}
			}

			w.Write([]byte("host not allowed"))
			w.WriteHeader(http.StatusUnauthorized)
		})
	}
}
