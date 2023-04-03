package certprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEmptyCertFileParam(t *testing.T) {
	params := []string{"--updated-proxy-listener-key-file=/certs/key-file.pem"}
	cp, err := new(Factory).New(params)
	assert.NotNil(t, err)
	assert.Nil(t, cp)
	assert.Equal(t, "parameter updated-proxy-listener-cert-file is required", err.Error())
}

func TestNewEmptyKeyFileParam(t *testing.T) {
	params := []string{"--updated-proxy-listener-cert-file=/certs/cert-file.pem"}
	cp, err := new(Factory).New(params)
	assert.NotNil(t, err)
	assert.Nil(t, cp)
	assert.Equal(t, "parameter updated-proxy-listener-key-file is required", err.Error())
}

func TestNewEmptyCheckIntervalParam(t *testing.T) {
	params := []string{"--updated-proxy-listener-cert-file=/certs/cert-file.pem", "--updated-proxy-listener-key-file=/certs/key-file.pem"}
	cp, err := new(Factory).New(params)
	assert.NotNil(t, err)
	assert.Nil(t, cp)
	assert.Equal(t, "parameter update-check-interval-minutes is required", err.Error())
}
