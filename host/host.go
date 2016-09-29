// Package host contains Middlewares for connection between
// hosts.
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
