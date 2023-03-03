package certprovider

import (
	"testing"
	"time"

	"github.com/grepplabs/kafka-proxy/pkg/libs/util"
	test_util "github.com/grepplabs/kafka-proxy/pkg/libs/util/test_util"
	"github.com/stretchr/testify/assert"
)

func TestGetX509KeyPair(t *testing.T) {
	testCertExpiryDate, dateErr := time.Parse(util.CertTimeStampLayout, test_util.TestCertExpiryDate)
	assert.Nil(t, dateErr, "failed to parse test date of certificate expiry")
	certificateProvider := &CertificateProvider{
		certPEMBlock:      []byte(test_util.TestCertPEM),
		keyPEMBlock:       []byte(test_util.TestKeyPEM),
		certificateExpiry: testCertExpiryDate,
	}

	resCertPEMBlock, resKeyPEMBlock, resErr := certificateProvider.GetX509KeyPair("2023-03-02T13:29:09Z")
	assert.Nil(t, resErr)
	assert.Equal(t, []byte(test_util.TestCertPEM), resCertPEMBlock)
	assert.Equal(t, []byte(test_util.TestKeyPEM), resKeyPEMBlock)
}

func TestGetX509KeyPairNoReturn(t *testing.T) {
	testCertExpiryDate, dateErr := time.Parse(util.CertTimeStampLayout, test_util.TestCertExpiryDate)
	assert.Nil(t, dateErr, "failed to parse test date of certificate expiry")
	certificateProvider := &CertificateProvider{
		certPEMBlock:      []byte(test_util.TestCertPEM),
		keyPEMBlock:       []byte(test_util.TestKeyPEM),
		certificateExpiry: testCertExpiryDate,
	}

	resCertPEMBlock, resKeyPEMBlock, resErr := certificateProvider.GetX509KeyPair(test_util.TestCertExpiryDate)
	assert.Nil(t, resErr)
	assert.Nil(t, resCertPEMBlock)
	assert.Nil(t, resKeyPEMBlock)
}
