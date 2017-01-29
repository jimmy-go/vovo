// Package thrott contains rate limiter handler.
//
// The MIT License (MIT)
//
// Copyright (c) 2016 Angel Del Castillo
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package thrott

import (
	"log"
	"net/http"
	"time"

	"github.com/didip/tollbooth"
)

// Limit limits the request by client to count by second.
func Limit(count int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		limiter := tollbooth.NewLimiter(count, time.Second)
		limiter.IPLookups = []string{"X-Real-IP", "RemoteAddr", "X-Forwarded-For"}

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tollbooth.SetResponseHeaders(limiter, w)

			httpError := tollbooth.LimitByRequest(limiter, r)
			if httpError != nil {
				w.Header().Add("Content-Type", limiter.MessageContentType)
				w.WriteHeader(httpError.StatusCode)
				_, err := w.Write([]byte(httpError.Message))
				if err != nil {
					log.Printf("write response err [%s]", err)
				}
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
