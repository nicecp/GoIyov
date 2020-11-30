package entity

import (
	"bufio"
	"bytes"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type Entity struct {
	startTime,endTime time.Time
	Request *http.Request
	Response *http.Response
	reqBody,  respBody   io.ReadCloser
}

func NewEntity(conn net.Conn) (*Entity, error) {
	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		return nil, errors.Wrap(err, "请求对象生成失败")
	}

	// http.Request.Body can only be read once, a new body needs to be copied
	body, err := getBody(request)
	if err != nil {
		return nil, err
	}

	request.URL.Scheme = "https"
	request.Body = body
	return &Entity{
		startTime: time.Now(),
		Request: request,
		reqBody: body,
	}, nil
}


func NewEntityWithRequest(request *http.Request) (*Entity, error) {
	body, err := getBody(request)
	if err != nil {
		return nil, err
	}

	request.Body = body
	return &Entity{
		startTime: time.Now(),

		Request: request,
		reqBody: body,
	}, nil
}

func (entity *Entity) setResponse(response *http.Response) {
	entity.endTime = time.Now()
}

func (entity *Entity) SetScheme(scheme string) *Entity {
	entity.Request.URL.Scheme = scheme
	return entity
}

func (entity *Entity) SetHost(host string) *Entity {
	entity.Request.URL.Host = host
	return entity
}

func (entity *Entity) SetRemoteAddr(remoteAddr string) *Entity {
	entity.Request.RemoteAddr = remoteAddr
	return entity
}

func getBody(request *http.Request) (io.ReadCloser, error) {
	reqBody, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, errors.Wrap(err, "获取请求Body失败")
	}
	return ioutil.NopCloser(bytes.NewReader(reqBody)), nil
}

