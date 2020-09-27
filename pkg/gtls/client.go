package gtls

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"google.golang.org/grpc/credentials"
)

type Client struct {
	ServerName string
	CaFile     string
	CertFile   string
	KeyFile    string
}

func (c *Client) GetCredentialsByCA() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(c.CaFile)
	if err != nil {
		return nil, err
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, err
	}

	client := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   c.ServerName,
		RootCAs:      certPool,
	})

	return client, err
}

func (c *Client) GetTLSCredentials() (credentials.TransportCredentials, error) {
	client, err := credentials.NewClientTLSFromFile(c.CertFile, c.ServerName)
	if err != nil {
		return nil, err
	}
	return client, err
}
