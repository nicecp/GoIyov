package proxy

import (
	"net/http"
	"time"
)

func DoRequest(req * http.Request) (*http.Response, error){
	transport := &http.Transport{
		DisableKeepAlives: true,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	req.RequestURI = ""
	return (&http.Client{Transport:transport}).Do(req)
}
