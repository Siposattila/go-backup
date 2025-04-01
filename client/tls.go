package client

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/Siposattila/gobkup/cert"
	"github.com/Siposattila/gobkup/log"
)

func (c *client) getTlsConfig() {
	ca, _, caError := cert.GenerateCA()
	if caError != nil {
		log.GetLogger().Fatal("Unable to generate CA certificate.", caError.Error())
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(ca)
	tlsConfig := &tls.Config{RootCAs: certPool, InsecureSkipVerify: true}
	c.Dialer.TLSClientConfig = tlsConfig
	log.GetLogger().Success("Tls config was obtained successfully!")
}
