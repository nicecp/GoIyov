package IyovGo

import (
	"IyovGo/cert"
	"IyovGo/conn"
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	tunnelConnectionEstablished = []byte("HTTP/1.1 200 Connection Established\r\n\r\n") // 通道连接建立
	badGateway = []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n\rn", http.StatusBadGateway, http.StatusText(http.StatusBadGateway)))
	hopToHopHeader = []string{ // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers
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

type Proxy struct {

}

func (proxy *Proxy)ServerHandler(rw http.ResponseWriter, req *http.Request) {
	clientConn, err := conn.HijackerConn(rw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	switch req.Method {
	case http.MethodConnect: // https
		clientConn.Write(tunnelConnectionEstablished)
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
		clientConn.Write(badGateway)
		return
	}

	tlsConn := tls.Server(clientConn,&tls.Config{Certificates: []tls.Certificate{certificate}})
	if err := tlsConn.Handshake(); err != nil {
		fmt.Printf("%+v",errors.WithStack(err))
	}

	tlsConn.SetDeadline(time.Now().Add(30 * time.Second))
	defer tlsConn.Close()
	request, err := http.ReadRequest(bufio.NewReader(tlsConn))
	if err != nil {
		fmt.Printf("%+v",errors.WithStack(err))
		clientConn.Write([]byte("TLS链接请求读取失败"))
	}
	request.URL.Scheme = "https"
	request.URL.Host = req.URL.Host

	// http.Request.Body can only be read once, a new body needs to be copied
	reqBody, err := ioutil.ReadAll(request.Body)
	request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

	resp, err := proxy.doRequest(request)
	if err != nil {
		fmt.Printf("%+v",errors.WithStack(err))
		clientConn.Write(badGateway)
		return
	}

	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewReader(respBody))
	resp.Write(tlsConn)
	//request.Body = ioutil.NopCloser(bytes.NewReader(reqBody))

	//a, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(respBody))

}

func (proxy *Proxy)handleHTTP(clientConn *conn.Connection, req *http.Request){
	defer clientConn.Close()
	resp, err := proxy.doRequest(req)
	if err != nil {
		clientConn.Write(badGateway)
		return
	}
	defer resp.Body.Close()
	resp.Write(clientConn)
}

// 请求目标服务器
func (proxy *Proxy)doRequest(req *http.Request) (*http.Response, error) {
	removeHopHeader(req.Header)
	resp, err :=  (&http.Transport{
		//DisableKeepAlives: true,
		ResponseHeaderTimeout: 30 * time.Second,
	}).RoundTrip(req)
	if err != nil {
		return nil, err
	}
	removeHopHeader(resp.Header)
	return resp, err
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