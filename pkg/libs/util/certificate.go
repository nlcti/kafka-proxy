package util

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"time"

	"github.com/pkg/errors"
)

const CertTimeStampLayout = time.RFC3339

// parseCertificate - parses given []byte into a certificate, if input is a valid certificate PEM block
func parseCertificate(certPEMBlock []byte) (*x509.Certificate, error) {

	block, _ := pem.Decode(certPEMBlock)

	cert, parseErr := x509.ParseCertificate(block.Bytes)

	if parseErr != nil {
		return nil, errors.Wrap(parseErr, "failed to parse certificate")
	}

	return cert, nil
}

// ExtractExpiryFromCertificate - reads certificate from given PEM block and extracts its expiry timestamp
func ExtractExpiryFromCertificate(certPEMBlock []byte) (time.Time, error) {
	certificate, err := parseCertificate(certPEMBlock)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to parse certificate")
	}

	return certificate.NotAfter, nil
}

// ExtractExpiryFromCertificateFile - reads certificate from given file path and extracts its expiry timestamp
func ExtractExpiryFromCertificateFile(certFile string) (time.Time, error) {
	certificate, err := ParseCertificateFile(certFile)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "failed to parse certificate")
	}

	return certificate.NotAfter, nil
}

// ParseCertificateFile - reads certificate from given file path and tries to parse its content
func ParseCertificateFile(certFile string) (*x509.Certificate, error) {

	content, readErr := ioutil.ReadFile(certFile)

	if readErr != nil {
		return nil, errors.Errorf("Failed to read file from location '%s'", certFile)
	}

	cert, parseErr := parseCertificate(content)
	if parseErr != nil {
		return nil, errors.Errorf("Failed to parse certificate file from location '%s'", certFile)
	}
	return cert, nil
}
