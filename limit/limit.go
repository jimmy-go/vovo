// Package limit contains tools to set hard-limit and soft-limit.
package limit

import (
	"net/http"
	"syscall"
)

// Hard set hard limit for application.
func Hard(x uint64) error {
	var rLimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err
	}

	rLimit.Max = x
	rLimit.Cur = x

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err
	}
	err = syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit)
	if err != nil {
		return err
	}

	return nil
}

// MaxBytes convenience for http.MaxBytesReader.
func MaxBytes(n int64) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			h.ServeHTTP(w, r)
		})
	}
}
