package main

import (
	"github.com/nicecp/GoIyov"
	"net/http"
	"time"
)

func main() {
	proxy := GoIyov.New()
	server := &http.Server{
		Addr: ":8888",
		Handler: http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			proxy.ServerHandler(rw, req)
		}),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
