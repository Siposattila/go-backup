package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"
)

func GenerateCA() (*x509.Certificate, *rsa.PrivateKey, error) {
	certTemplate := &x509.Certificate{
		SerialNumber:          big.NewInt(2024),
		Subject:               pkix.Name{},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	caPrivateKey, generationError := rsa.GenerateKey(rand.Reader, 2048)
	if generationError != nil {
		return nil, nil, generationError
	}

	caBytes, creationErr := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, &caPrivateKey.PublicKey, caPrivateKey)
	if creationErr != nil {
		return nil, nil, creationErr
	}

	ca, parseError := x509.ParseCertificate(caBytes)
	if parseError != nil {
		return nil, nil, parseError
	}

	return ca, caPrivateKey, nil
}

func GenerateLeafCert(dnsName string, ca *x509.Certificate, caPrivateKey *rsa.PrivateKey) (*x509.Certificate, *rsa.PrivateKey, error) {
	certTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		DNSNames:     []string{dnsName},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	privKey, genrationError := rsa.GenerateKey(rand.Reader, 2048)
	if genrationError != nil {
		return nil, nil, genrationError
	}

	certBytes, creationError := x509.CreateCertificate(rand.Reader, certTemplate, ca, &privKey.PublicKey, caPrivateKey)
	if creationError != nil {
		return nil, nil, creationError
	}

	cert, parseError := x509.ParseCertificate(certBytes)
	if parseError != nil {
		return nil, nil, parseError
	}

	return cert, privKey, nil
}
