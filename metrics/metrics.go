// Package metrics contains prometheus metrics for requests,
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
package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(Clients)
	prometheus.MustRegister(Durations)
	prometheus.MustRegister(Errors)
}

var (
	// Clients register requests count.
	Clients = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "proy",
			Name:      "requests",
			Help:      "Number of petitions.",
		},
		[]string{"Clients"},
	)
	// Durations duration of every request.
	Durations = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "proy",
			Name:      "duration_miliseconds",
			Help:      "Duration in miliseconds.",
		},
		[]string{"Duration"},
	)
	// Errors register errors count.
	Errors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "proy",
			Name:      "errors",
			Help:      "Errors count.",
		},
		[]string{"Errors"},
	)
)

// ClientInc increments clients petitions by label.
func ClientInc(s string) {
	Clients.WithLabelValues(s).Inc()
}

// ErrorInc increments errors by label.
func ErrorInc(s string) {
	Errors.WithLabelValues(s).Inc()
}

// DurationObs register petition time by label.
func DurationObs(start time.Time, s string) {
	tot := time.Since(start)
	Durations.WithLabelValues(s).Observe(float64(tot) / 1000000)
}

// Handler metrics middleware.
func Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// skip prometheus default endpoint
		if r.RequestURI == "/metrics" {
			return
		}

		start := time.Now()
		h.ServeHTTP(w, r)
		// slug := fmt.Sprintf("%v-%v", r.URL.Path, r.Method)
		slug := r.URL.Path + "-" + r.Method
		DurationObs(start, slug)
	})
}

// Custom metrics middleware.
// use this when you want a control over paths with params
// like /me/:var1/size/:var2 convert it to: /me/sizes
//
func Custom(path string) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// skip prometheus default endpoint
			if r.RequestURI == "/metrics" {
				return
			}

			start := time.Now()
			h.ServeHTTP(w, r)
			DurationObs(start, path)
		})
	}
}
