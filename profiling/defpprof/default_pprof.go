// Package defpprof contains a server for net profiling.
/*
	Usage:

		_ "github.com/jimmy-go/profiling/defpprof"

	will start a server on localhost:6060
*/
package defpprof

import (
	"log"
	"net/http"
	"runtime"

	// import pprof
	_ "net/http/pprof"
)

func init() {
	log.Printf("Multi init")
	runtime.SetBlockProfileRate(1)
	go pprofServer()
}

func pprofServer() {
	log.Printf("Multi : web pprofiling enabled : listening :6060")
	log.Println(http.ListenAndServe("localhost:6060", nil))
}
