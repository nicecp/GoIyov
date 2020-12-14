package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/nicecp/GoIyov/cache"
	"github.com/pkg/errors"
	"math/big"
	"net"
	"time"
)

var (
	rootCa  *x509.Certificate // CA证书
	rootKey *rsa.PrivateKey   // 证书私钥
)

var (
	_rootCa = []byte(`-----BEGIN CERTIFICATE-----
MIIDmjCCAoICCQCe/26mrL7IqzANBgkqhkiG9w0BAQUFADCBjjELMAkGA1UEBhMC
Q04xEDAOBgNVBAgMB1NpY2h1YW4xEDAOBgNVBAcMB0NoZW5nZHUxDzANBgNVBAoM
Bkl5b3ZHbzEPMA0GA1UECwwGSXlvdkdvMQ8wDQYDVQQDDAZDYVJvb3QxKDAmBgkq
hkiG9w0BCQEWGWZvcndhcmQubmljZS5jcEBnbWFpbC5jb20wHhcNMjAxMTI0MDI1
NTMxWhcNMjExMTI0MDI1NTMxWjCBjjELMAkGA1UEBhMCQ04xEDAOBgNVBAgMB1Np
Y2h1YW4xEDAOBgNVBAcMB0NoZW5nZHUxDzANBgNVBAoMBkl5b3ZHbzEPMA0GA1UE
CwwGSXlvdkdvMQ8wDQYDVQQDDAZDYVJvb3QxKDAmBgkqhkiG9w0BCQEWGWZvcndh
cmQubmljZS5jcEBnbWFpbC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEK
AoIBAQC0HyQBJLLly7+/W879qBSTHycEQNO87G89L3qIaTOkWCjwPRfPh9Xgkaie
o6g8iJVB023ytUqPOMjRJ2H4kWIZIRNsoWhkdlfw2WaFcLoA9gE+MCKFN/6xmN+j
b7dSnF/amdw6tuxrJl2/1vFx4IxRttACfjnDO1IluY/zaPDy4zhxZIWW2o5kKACH
CPg1bqN8M5ifprfvEvwEJ8e9CmtLVHXCfg83hXvRrrGgtL+YBW89xjZYMxjEN8xT
6bmji8QNu3i/NrJhwQIX15WJCpFOKkGlKzScML4+yzqLtJYeKVGdRciEnSBZT4dG
r2lJKN/ZlGTVTXm5FL7NN41U2OgbAgMBAAEwDQYJKoZIhvcNAQEFBQADggEBAAmk
ITfrki7GMqu+OJuHAlV6YGaWd3r8Y+DGgYFYsIP36OkQWbdjwfkx0PiwYgWITUWE
RKrNWNMvT43FuFPLaaQ/i7Of84+QhcB2oCIe1exv/cq0kS/b4pM3qBSufEFZsb6O
me+tQLS7jPeA/D9GGGjwl20KDv2Q8bqbog+jiV6JygF74r3ByAmaECLFFOmBCRzL
Wp6GRz7uTmx42hLsGSJwMSPEUzjo23WbICZblGmwFFsep/3Ly9yuNVgVQgsZD5A/
0WN8PhtnroGQ45TQa5ObbiJQzru/84xe8YTUdn0B/4k+af8C0qdjTPSypexBxqEq
uJ0ddXfqds6e3Z8ZUaU=
-----END CERTIFICATE-----
`)
	_rootKey = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIIEpgIBAAKCAQEAtB8kASSy5cu/v1vO/agUkx8nBEDTvOxvPS96iGkzpFgo8D0X
z4fV4JGonqOoPIiVQdNt8rVKjzjI0Sdh+JFiGSETbKFoZHZX8NlmhXC6APYBPjAi
hTf+sZjfo2+3Upxf2pncOrbsayZdv9bxceCMUbbQAn45wztSJbmP82jw8uM4cWSF
ltqOZCgAhwj4NW6jfDOYn6a37xL8BCfHvQprS1R1wn4PN4V70a6xoLS/mAVvPcY2
WDMYxDfMU+m5o4vEDbt4vzayYcECF9eViQqRTipBpSs0nDC+Pss6i7SWHilRnUXI
hJ0gWU+HRq9pSSjf2ZRk1U15uRS+zTeNVNjoGwIDAQABAoIBAQCYbnoj1ZDoVAOT
x/hmReYTk5uLR+loypZhK1sBMjaX8FvE447Q/F2NzPbsOgfYIqZdrLYxXicZCa85
AaExoKdqKMmtdvNHgbduhizy5LEkuwvWOxobr4WFeqBYSeTUrq2X6/mqXr+49iEE
hryR6LwXMyTZ10S+6ebdMiqWjcrLYP3dVGcmD29hcBgPutOh1g5jI41L5qKS6PJ5
qqRkGsZpYWyonNLbGdcVeZOfguENHUfZyoEjt0z+EXlsd6zYczLFZOOiyDDLg+vX
8xcEtKbhD4s3xnaxPuZZQOoNms6uZPd0dXGqknastTlHi8SfCz1jw7lvPJYuH4hM
GDHDTSGhAoGBAPAV/ixvAgx8PJRzcqLw4oznttodgaAgFHQ3Jl0b8GU/NPV5P+Uv
CRuhFTS4+hFcS3kCItmzgHbQlcGbchOknd52SdLDJrzj3/L6k+4/jeIEcx1MaQYc
8YyRrEB8+5myK/BvQ986xWj+emTA4Xgpy5pgEba4R1tTtM9smBuF62S5AoGBAMAP
nebIrM7Rs6NkizZGP0yWmgGT/WLMa74UV32+kDjR9LG6521Sfe9dc91CUOmEsOa4
EYk9bCWBDt/6kcSDxWIYPnkhS0FPz3d34tQbbrENcNCVMuyhnibPXTzO9NLWYC9/
iUzxYqeLhSvXza7A0MPnv5pRBm0SvDrpYstBTXFzAoGBAMDfWi8GCuZO1DgKOvjt
fYLnD31QIPerbeMi/v3j2Q7tZTUi8BLE45M/qBKP280gkT0oWyj7TGOnE/fSUiW3
pF+4NXxM7Izon9vKNBc9FVWSb4wE+4Y+sEpWKMQx48pIWYYxTJxD0Z2Uem0AiuGG
6hsdvH1Gs4SJzYKpYdUSk9V5AoGBAJ2HsaHrkyIICmnIPA8GS0EMfcExmzGALhc4
JBL1TOHuA+ALR2r5sGW2pyQiEq+WsGptK6T/hka0tnir0wf2dN1iuUstLcaiKa75
3EjRP1dliNTsq1o/rbJzfywzK8gLIdWTrBA6JQr7ev1dAk2FxTYKTbPLJZQtO8qu
RuQj6dtVAoGBAMlYAQncj8fF5MQORf9JuAv1OacIBhzorEaIBjIsD713sH8nXrkE
IYtgRQrk5EBkRVPrR3CVTvHo4fqFcUIOgViy/03azoslVZ5ZvQtBBTzp8mq98H8E
jazsVKrVKRxltQUJ0m+AZSfL5s53fz6L1qADE2Vzgqibc8t7CPrDEOgu
-----END RSA PRIVATE KEY-----
`)
)

var certCache *cache.Cache

func init() {
	certCache = cache.NewCache()

	if err := loadRootCa(); err != nil {
		panic(err)
	}
	if err := loadRootKey(); err != nil {
		panic(err)
	}
}

func GetCertificate(host string) (tls.Certificate, error) {
	certificate, err := certCache.GetOrStore(host, func() (interface{}, error) {
		host, _, err := net.SplitHostPort(host)
		if err != nil {
			return nil, err
		}
		certByte, priByte, err := generatePem(host)
		if err != nil {
			return nil, err
		}
		certificate, err := tls.X509KeyPair(certByte, priByte)
		if err != nil {
			return nil, err
		}
		return certificate, nil
	})
	return certificate.(tls.Certificate), err
}
func generatePem(host string) ([]byte, []byte, error) {
	max := new(big.Int).Lsh(big.NewInt(1), 128)   //把 1 左移 128 位，返回给 big.Int
	serialNumber, _ := rand.Int(rand.Reader, max) //返回在 [0, max) 区间均匀随机分布的一个随机值
	template := x509.Certificate{
		SerialNumber: serialNumber, // SerialNumber 是 CA 颁布的唯一序列号，在此使用一个大随机数来代表它
		Subject: pkix.Name{ //Name代表一个X.509识别名。只包含识别名的公共属性，额外的属性被忽略。
			CommonName: host,
		},
		NotBefore:      time.Now().AddDate(-1, 0, 0),
		NotAfter:       time.Now().AddDate(1, 0, 0),
		KeyUsage:       x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, //KeyUsage 与 ExtKeyUsage 用来表明该证书是用来做服务器认证的
		ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},               // 密钥扩展用途的序列
		EmailAddresses: []string{"forward.nice.cp@gmail.com"},
	}

	if ip := net.ParseIP(host); ip != nil {
		template.IPAddresses = []net.IP{ip}
	} else {
		template.DNSNames = []string{host}
	}

	priKey, err := generateKeyPair()
	if err != nil {
		return nil, nil, err
	}

	cer, err := x509.CreateCertificate(rand.Reader, &template, rootCa, &priKey.PublicKey, rootKey)
	if err != nil {
		return nil, nil, err
	}

	return pem.EncodeToMemory(&pem.Block{ // 证书
			Type:  "CERTIFICATE",
			Bytes: cer,
		}), pem.EncodeToMemory(&pem.Block{ // 私钥
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priKey),
		}), err
}

// 秘钥对 生成一对具有指定字位数的RSA密钥
func generateKeyPair() (*rsa.PrivateKey, error) {
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.Wrap(err, "密钥对生成失败")
	}

	return priKey, nil
}

// 加载根证书
func loadRootCa() error {
	p, _ := pem.Decode(_rootCa)
	var err error
	rootCa, err = x509.ParseCertificate(p.Bytes)
	if err != nil {
		return errors.Wrap(err, "CA证书解析失败")
	}

	return nil
}

// 加载根Private Key
func loadRootKey() error {
	p, _ := pem.Decode(_rootKey)
	var err error
	rootKey, err = x509.ParsePKCS1PrivateKey(p.Bytes)
	if err != nil {
		return errors.Wrap(err, "Key证书解析失败")
	}

	return err
}
