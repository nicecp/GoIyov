package proxy

import (
	"IyovGo/conn"
	"fmt"
	"net/http"
	"time"
)

type Proxy struct {

}

func New () *Proxy {
	return &Proxy{}
}
// 通道连接建立
var tunnelConnectionEstablished = []byte("HTTP/1.1 200 Connection Established\r\n\r\n")
var badGateway = []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n\rn", http.StatusBadGateway, http.StatusText(http.StatusBadGateway)))

func (proxy *Proxy)ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	clientConn, err := conn.HijackerConn(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer clientConn.Close()
	switch req.Method {
	case http.MethodConnect: // https
		_,_ = clientConn.Write(badGateway)
		//_,err = clientConn.Write([]byte("暂不支持https"))
	default :
		// http
		// todo websocket
		resp, err := handleHTTP(req)
		if err != nil {
			return
		}
		defer resp.Body.Close()
		resp.Write(clientConn)
	}
}

func handleHTTP(req * http.Request) (*http.Response, error){
	//req.RequestURI = ""
	//return (&http.Client{Transport:transport}).Do(req)
	return (&http.Transport{
		DisableKeepAlives: true,
		ResponseHeaderTimeout: 30 * time.Second,
	}).RoundTrip(req)
}