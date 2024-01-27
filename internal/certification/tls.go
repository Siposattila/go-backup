package certification

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/Siposattila/gobkup/internal/config"
	"github.com/Siposattila/gobkup/internal/console"
)

const alpn = "webtransport-go / quic-go"

var TlsConfig  *tls.Config
var CertPool *x509.CertPool

func GetClientTlsConfig() {
	var ca, _, caError = generateCA()
	if caError != nil {
		console.Fatal("Unable to generate CA certificate: " + caError.Error())
	}

	CertPool = x509.NewCertPool()
	CertPool.AddCert(ca)
	TlsConfig = &tls.Config{RootCAs: CertPool}
}

func GetServerTlsConfig() {
	var ca, caPrivateKey, caError = generateCA()
	if caError != nil {
		console.Fatal("Unable to generate CA certificate: " + caError.Error())
	}

	var leafCert, leafPrivateKey, leafError = generateLeafCert(ca, caPrivateKey)
	if leafError != nil {
		console.Fatal("Unable to generate leaf certificate: " + leafError.Error())
	}

	CertPool = x509.NewCertPool()
	CertPool.AddCert(ca)
	TlsConfig = &tls.Config{
		Certificates: []tls.Certificate{{
			Certificate: [][]byte{leafCert.Raw},
			PrivateKey:  leafPrivateKey,
		}},
		NextProtos: []string{alpn},
	}
}

func generateCA() (*x509.Certificate, *rsa.PrivateKey, error) {
	var certTemplate = &x509.Certificate{
		SerialNumber:          big.NewInt(2024),
		Subject:               pkix.Name{},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	var caPrivateKey, generationError = rsa.GenerateKey(rand.Reader, 2048)
	if generationError != nil {
		return nil, nil, generationError
	}

	var caBytes, creationErr = x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if creationErr != nil {
		return nil, nil, creationErr
	}

	var ca, parseError = x509.ParseCertificate(caBytes)
	if parseError != nil {
		return nil, nil, parseError
	}

	return ca, caPrivateKey, nil
}

func generateLeafCert(ca *x509.Certificate, caPrivateKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey, error) {
	var certTemplate = &x509.Certificate{
		SerialNumber: big.NewInt(1),
		DNSNames:     []string{config.Master.Domain},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	var privKey, genrationError = rsa.GenerateKey(rand.Reader, 2048)
	if genrationError != nil {
		return nil, nil, genrationError
	}

	var certBytes, creationError = x509.CreateCertificate(rand.Reader, certTemplate, ca, &privKey.PublicKey, caPrivateKey)
	if creationError != nil {
		return nil, nil, creationError
	}

	var cert, parseError = x509.ParseCertificate(certBytes)
	if parseError != nil {
		return nil, nil, parseError
	}

	return cert, privKey, nil
}
