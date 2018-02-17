// Package profiling contains a server for net profiling.
package profiling

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	// import pprof
	_ "net/http/pprof"
)

// Listen will start a server on port.
func Listen(port int) {
	runtime.SetBlockProfileRate(1)
	go func(p int) {
		log.Printf("Multi : web pprofiling enabled : listening :%v", p)
		log.Println(http.ListenAndServe(fmt.Sprintf(":%v", p), nil))
	}(port)
}
