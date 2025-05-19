package server

import (
	"crypto/tls"

	"github.com/Siposattila/go-backup/cert"
	"github.com/Siposattila/go-backup/log"
)

func (s *server) getTlsConfig() {
	ca, caPrivateKey, caError := cert.GenerateCA()
	if caError != nil {
		log.GetLogger().Fatal("Unable to generate CA certificate.", caError.Error())
	}

	leafCert, leafPrivateKey, leafError := cert.GenerateLeafCert(s.Config.Domain, ca, caPrivateKey)
	if leafError != nil {
		log.GetLogger().Fatal("Unable to generate leaf certificate.", leafError.Error())
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{leafCert.Raw},
			PrivateKey:  leafPrivateKey,
		}},
		NextProtos:         []string{"webtransport-go / quic-go"},
		InsecureSkipVerify: true,
	}

	s.Transport.H3.TLSConfig = tlsConfig
	log.GetLogger().Success("Tls config was obtained successfully!")
}
