package util

import (
	_ "embed"
	"testing"

	test_util "github.com/grepplabs/kafka-proxy/pkg/libs/util/test_util"
	"github.com/stretchr/testify/assert"
)

func TestExtractExpiryFromCertificate(t *testing.T) {

	expiry, err := ExtractExpiryFromCertificate([]byte(test_util.TestCertPEM))
	assert.Nil(t, err)
	assert.False(t, expiry.IsZero())
	assert.Equal(t, test_util.TestCertExpiryDate, expiry.Format(CertTimeStampLayout))
}
