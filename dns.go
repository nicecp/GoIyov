package GoIyov

import (
	"net"
)

type Dns struct {
	SslCertHost string
	records map[string]string
}

var DefaultDns = Dns{
	SslCertHost: "goiyov",
	records: make(map[string]string),
}

// 添加DNS映射
func (dns *Dns) Add(records map[string]string) {
	for host, remote := range records {
		dns.records[host] = remote
	}
}

func (dns *Dns) CustomDialer(addr string) (string, error) {
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr, err
	}

	if destHost, ok := dns.records[host]; ok {
		return destHost + ":" + port, nil
	}
	return addr, nil
}
