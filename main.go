package main

import (
	"IyovGo/proxy"
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         ":8888",
		Handler:      proxy.New(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
