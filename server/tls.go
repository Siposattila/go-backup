package server

import (
	"crypto/tls"
	"log"

	"github.com/Siposattila/gobkup/cert"
)

func (s *server) getTlsConfig() {
	ca, caPrivateKey, caError := cert.GenerateCA()
	if caError != nil {
		log.Fatal("Unable to generate CA certificate: " + caError.Error())
	}

	leafCert, leafPrivateKey, leafError := cert.GenerateLeafCert("server", ca, caPrivateKey)
	if leafError != nil {
		log.Fatal("Unable to generate leaf certificate: " + leafError.Error())
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
	log.Println("Tls config was obtained successfully!")
}
