// Package cors contains middleware to enable site CORS.
// request duration and errors count.
// Middleware func allows easy handler register.
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
package cors

import (
	"net/http"
	"strings"
)

// New returns a middleware with allow-access-control-origins header.
// origins is a string with hosts separated by comma.
func New(origins string) func(http.Handler) http.Handler {
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

			w.Write([]byte("host not supported"))
		})
	}
}
