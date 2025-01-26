package server

import (
	"crypto/tls"

	"github.com/Siposattila/gobkup/cert"
	"github.com/Siposattila/gobkup/log"
)

func (s *server) getTlsConfig() {
	ca, caPrivateKey, caError := cert.GenerateCA()
	if caError != nil {
		log.GetLogger().Fatal("Unable to generate CA certificate: ", caError)
	}

	leafCert, leafPrivateKey, leafError := cert.GenerateLeafCert("server", ca, caPrivateKey)
	if leafError != nil {
		log.GetLogger().Fatal("Unable to generate leaf certificate: ", leafError)
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{leafCert.Raw},
			PrivateKey:  leafPrivateKey,
		}},
		NextProtos: []string{"webtransport-go / quic-go"},
	}

	tlsConfig.InsecureSkipVerify = true
	s.Transport.H3.TLSConfig = tlsConfig
	log.GetLogger().Success("Tls config was obtained successfully!")
}
