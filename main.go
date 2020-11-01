package IyovGo

import (
	"net/http"
	"time"
)

func main() {
	server := &http.Server{
		Addr:         ":8888",
		Handler:      proxy,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
