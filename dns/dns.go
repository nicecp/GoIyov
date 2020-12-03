package dns

import (
	"net"
)

var dnsMapping  = map[string]string{
	"localhost-x": "127.0.0.1",
}

func CustomDialer(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr, err
	}

	if destHost, ok := dnsMapping[host]; ok {
		return destHost + ":" + port, nil
	}
	return addr, nil
}