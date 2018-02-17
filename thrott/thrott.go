// Package thrott contains rate limiter handler.
package thrott

import (
	"net/http"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
)

// Limit limits the request by client to count by second.
func Limit(count float64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		lmt := tollbooth.NewLimiter(count, &limiter.ExpirableOptions{DefaultExpirationTTL: time.Minute})
		lmt.SetIPLookups([]string{"X-Real-IP", "RemoteAddr", "X-Forwarded-For"})

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpError := tollbooth.LimitByRequest(lmt, w, r)
			if httpError != nil {
				http.Error(w, httpError.Message, httpError.StatusCode)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
