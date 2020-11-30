package IyovGo

import (
	"IyovGo/cert"
	"IyovGo/conn"
	"IyovGo/entity"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	tunnelConnectionEstablished = []byte("HTTP/1.1 200 Connection Established\r\n\r\n") // 通道连接建立
	internalServerErr			= "HTTP/1.1 %d %s\r\n\r\n"
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
		//resp.Write(clientConn)
		//copyHeader(rw.Header(), resp.Header)
		//rw.WriteHeader(resp.StatusCode)
		//io.Copy(rw, resp.Body) // Header也要发给rw
	}
}

func (proxy *Proxy)handleHTTPS(clientConn *conn.Connection,req *http.Request)  {
	defer clientConn.Close()
	certificate, err := cert.GetCertificate(req.URL.Host)
	if err != nil {
		fmt.Printf("%+v",errors.WithStack(err))
		Error(clientConn, http.StatusInternalServerError, err)
		return
	}

	tlsConn := tls.Server(clientConn,&tls.Config{Certificates: []tls.Certificate{certificate}})
	if err := tlsConn.Handshake(); err != nil {
		fmt.Printf("%+v",errors.WithStack(err))
		return
	}

	tlsConn.SetDeadline(time.Now().Add(30 * time.Second))
	defer tlsConn.Close()

	proxyEntity, err := entity.NewEntity(tlsConn)
	if err != nil {
		Error(tlsConn, http.StatusInternalServerError, err)
		return
	}
	proxyEntity.SetHost(req.URL.Host).SetRemoteAddr(req.RemoteAddr)

	resp, err := proxy.doRequest(tlsConn, proxyEntity)
	if err != nil {
		Error(tlsConn, http.StatusInternalServerError, err)
		return
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewReader(respBody))
	resp.Write(tlsConn)
}

func (proxy *Proxy)handleHTTP(clientConn *conn.Connection, req *http.Request){
	defer clientConn.Close()

	proxyEntity, err := entity.NewEntityWithRequest(req)
	if err != nil {
		fmt.Printf("%+v", errors.WithStack(err))
		Error(clientConn, http.StatusInternalServerError, err)
		return
	}
	resp, err := proxy.doRequest(clientConn, proxyEntity)
	if err != nil {
		fmt.Printf("%+v", errors.WithStack(err))
		Error(clientConn, http.StatusInternalServerError, err)
		return
	}
	defer resp.Body.Close()

	resp.Write(clientConn)

	a, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(a))
}

// 请求目标服务器
func (proxy *Proxy)doRequest(clientConn net.Conn, entity *entity.Entity) (*http.Response, error) {
	removeHopHeader(entity.Request.Header)
	resp, err :=  (&http.Transport{
		DisableKeepAlives: true,
		ResponseHeaderTimeout: 30 * time.Second,
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

func Error(net net.Conn, code int, error error) {
	_, _ = net.Write([]byte(fmt.Sprintf(internalServerErr, code, http.StatusText(code))))
	if error != nil {
		_, _ = net.Write([]byte(error.Error()))
	}
}