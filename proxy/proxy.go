package proxy

import (
	"IyovGo/conn"
	"net/http"
)

func ServerHTTP(rw http.ResponseWriter, req *http.Request) {
	clientConn, err := conn.HijackerConn(rw)
	defer clientConn.Close()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	switch req.Method {
	case http.MethodConnect:
		rw.WriteHeader(http.StatusBadGateway)
		_,_ = rw.Write([]byte("暂不支持https"))
	case http.MethodGet:

	}
}