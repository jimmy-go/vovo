// Package metrics contains prometheus metrics for requests, request duration
// and errors count.
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
func Custom(path string) func(http.Handler) http.Handler {
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
