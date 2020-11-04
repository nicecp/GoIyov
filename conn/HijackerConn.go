package conn

import (
	"net"
	"net/http"
)

type Connection struct {
	net.Conn
}

func HijackerConn(rw http.ResponseWriter) (*Connection,error) {
	hijacker, ok := rw.(http.Hijacker)
	if !ok {
		return nil, http.ErrHijacked
	}

	conn, _, err := hijacker.Hijack()
	if err != nil {
		return nil, err
	}

	return &Connection{conn}, nil
}
