package proxy

import (
	"crypto/tls"
	"errors"
	"testing"

	"github.com/grepplabs/kafka-proxy/pkg/libs/util"
	"github.com/grepplabs/kafka-proxy/pkg/libs/util/test_util"
	"github.com/stretchr/testify/assert"
)

type certificateProviderMock struct {
	testCertPEM []byte
	testKeyPEM  []byte
	testError   error
}

func (cpm *certificateProviderMock) GetX509KeyPair(expiryDate string) ([]byte, []byte, error) {
	return cpm.testCertPEM, cpm.testKeyPEM, cpm.testError
}

func TestGetCertificate(t *testing.T) {
	helloInfo := &tls.ClientHelloInfo{}
	cp := &certificateProviderMock{
		testCertPEM: []byte(test_util.TestCertPEM),
		testKeyPEM:  []byte(test_util.TestKeyPEM),
	}

	certificateClient := NewCertificateClient(cp)
	cert, err := certificateClient.GetCertificate(helloInfo)
	assert.Nil(t, err)
	assert.NotNil(t, cert)
	assert.Equal(t, test_util.TestCertExpiryDate, certificateClient.certificateExpiry.Format(util.CertTimeStampLayout))
	assert.True(t, certificateClient.certificateAvailable)
}

func TestGetCertificateWithError(t *testing.T) {
	helloInfo := &tls.ClientHelloInfo{}
	cp := &certificateProviderMock{
		testError: errors.New("TestError"),
	}

	certificateClient := NewCertificateClient(cp)
	cert, err := certificateClient.GetCertificate(helloInfo)
	assert.Nil(t, cert)
	assert.NotNil(t, err)
	assert.Equal(t, "failed to request X509 key pair from certificate provider plugin: TestError", err.Error())
}

func TestGetCertificateWithMissingKeyPEM(t *testing.T) {
	helloInfo := &tls.ClientHelloInfo{}
	cp := &certificateProviderMock{
		testCertPEM: []byte(test_util.TestCertPEM),
	}

	certificateClient := NewCertificateClient(cp)
	cert, err := certificateClient.GetCertificate(helloInfo)
	assert.Nil(t, cert)
	assert.NotNil(t, err)
	assert.Equal(t, "no TLS certificate available", err.Error())
}
