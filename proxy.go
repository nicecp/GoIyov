package IyovGo

import (
	"IyovGo/cert"
	"IyovGo/conn"
	"IyovGo/entity"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	tunnelConnectionEstablished = []byte("HTTP/1.1 200 Connection Established\r\n\r\n") // 通道连接建立
	internalServerErr			= []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)))
	hopToHopHeader              = []string{ // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers
		"Keep-Alive",
		"Transfer-Encoding",
		"TE",
		"Connection",
		"Trailer",
		"Upgrade",
		"Proxy-Authorization",
		"Proxy-Authenticate",
		"Connection",
	}
)

type Proxy struct {}

func (proxy *Proxy)ServerHandler(rw http.ResponseWriter, req *http.Request) {
	clientConn, err := conn.HijackerConn(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch req.Method {
	case http.MethodConnect: // https
		_, _ = clientConn.Write(tunnelConnectionEstablished)

		go proxy.handleHTTPS(clientConn, req)
	default : // todo websocket
		proxy.handleHTTP(clientConn, req)
	}
}

func (proxy *Proxy)handleHTTPS(clientConn *conn.Connection,req *http.Request)  {
	defer clientConn.Close()
	certificate, err := cert.GetCertificate(req.URL.Host)
	if err != nil {
		fmt.Printf("%+v",errors.WithStack(err))
		Error(clientConn, err)
		return
	}

	tlsConn := tls.Server(clientConn,&tls.Config{Certificates: []tls.Certificate{certificate}})
	if err := tlsConn.Handshake(); err != nil {
		fmt.Printf("%+v",errors.WithStack(err))
		return
	}

	_ = tlsConn.SetDeadline(time.Now().Add(30 * time.Second))
	defer tlsConn.Close()

	proxyEntity,err := entity.NewEntity(tlsConn)
	if err != nil {
		Error(tlsConn, err)
		return
	}

	proxyEntity.SetScheme("https")
	proxyEntity.SetHost(req.URL.Host)
	proxyEntity.SetRemoteAddr(req.RemoteAddr)

	resp, err := proxy.doRequest(tlsConn, proxyEntity)
	if err != nil {
		Error(tlsConn, err)
		return
	}

	defer resp.Body.Close()

	if err = proxyEntity.SetResponse(resp); err != nil {
		Error(tlsConn, err)
	}

	_ = resp.Write(tlsConn)
}

func (proxy *Proxy)handleHTTP(clientConn *conn.Connection, req *http.Request){
	defer clientConn.Close()

	proxyEntity, err := entity.NewEntityWithRequest(req)
	if err != nil {
		fmt.Printf("%+v", errors.WithStack(err))
		Error(clientConn, err)
		return
	}
	resp, err := proxy.doRequest(clientConn, proxyEntity)
	if err != nil {
		fmt.Printf("%+v", errors.WithStack(err))
		Error(clientConn, err)
		return
	}
	defer resp.Body.Close()

	_ = resp.Write(clientConn)

}

// 请求目标服务器
func (proxy *Proxy)doRequest(clientConn net.Conn, entity *entity.Entity) (*http.Response, error) {
	removeHopHeader(entity.Request.Header)

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
		Deadline: time.Now().Add(30 * time.Second),
	}
	resp, err :=  (&http.Transport{
		DisableKeepAlives: true,
		ResponseHeaderTimeout: 30 * time.Second,
		DialContext: func(ctx context.Context, network, addr string) (i net.Conn, e error) {
			addr, _ = CustomDialer(addr)
			return dialer.DialContext(ctx, network, addr)
		},
	}).RoundTrip(entity.Request)
	if err != nil {
		return nil, err
	}
	removeHopHeader(resp.Header)
	return resp, nil
}

// remove hop header
func removeHopHeader(header http.Header) {
	for _, hop := range hopToHopHeader {
		if value := header.Get(hop); len(value) != 0 {
			if strings.EqualFold(hop, "Connection") {
				for _, customerHeader := range strings.Split(value, ",") {
					header.Del(strings.Trim(customerHeader, " "))
				}
			}
			header.Del(hop)
		}
	}
}

func Error(net net.Conn, error error) {
	_, _ = net.Write(internalServerErr)
	if error != nil {
		_, _ = net.Write([]byte(error.Error()))
	}
}