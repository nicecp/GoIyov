package conn

import (
	"net"
	"net/http"
)

type Connection struct {
	conn net.Conn
}

func (connection *Connection) Close() {
	_ = connection.conn.Close()
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

func (connection *Connection) Write(b []byte) (int, error)  {
	return connection.conn.Write(b)
}
