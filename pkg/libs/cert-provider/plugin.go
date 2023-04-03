package certprovider

import (
	"crypto/tls"
	"os"
	"sync"
	"time"

	"github.com/grepplabs/kafka-proxy/pkg/libs/util"
	"github.com/pkg/errors"
)

type CertificateProvider struct {
	mu                sync.RWMutex
	certFile          string
	keyFile           string
	checkIntervalMin  int
	certPEMBlock      []byte
	keyPEMBlock       []byte
	certificateExpiry time.Time
}

type CertificateProviderOptions struct {
	CertFile         string
	KeyFile          string
	CheckIntervalMin int
}

// GetX509KeyPair - getter function to return
// PEM encoded blocks of actual certificate and key
func (cp *CertificateProvider) GetX509KeyPair(expiryDate string) ([]byte, []byte, error) {
	givenCertificateExpiry, err := time.Parse(util.CertTimeStampLayout, expiryDate)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to parse request parameter into certificate expiry")
	}
	// return data on new certificate
	if !givenCertificateExpiry.IsZero() && cp.certificateExpiry.After(givenCertificateExpiry) {
		cp.mu.RLock()
		defer cp.mu.RUnlock()
		return cp.certPEMBlock, cp.keyPEMBlock, nil
	}

	// no NEW certificate data
	return nil, nil, nil
}

// NewCertificateProvider - creates a new CertificateProvider from given specs
func NewCertificateProvider(opts CertificateProviderOptions) (*CertificateProvider, error) {

	if opts.CertFile == "" {
		return nil, errors.New("parameter updated-proxy-listener-cert-file is required")
	}
	if opts.KeyFile == "" {
		return nil, errors.New("parameter updated-proxy-listener-key-file is required")
	}
	if opts.CheckIntervalMin == 0 {
		return nil, errors.New("parameter update-check-interval-minutes is required")
	}

	certificateProvider := &CertificateProvider{certFile: opts.CertFile, keyFile: opts.KeyFile, checkIntervalMin: opts.CheckIntervalMin}

	// create certificate initially
	certExpiry, err := util.ExtractExpiryFromCertificateFile(certificateProvider.certFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to extract certificate expiry initially")
	}
	if certExpiry.IsZero() {
		return nil, errors.New("extracted invalid certificate expiry initially")
	}
	err = certificateProvider.extractNewCertificate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create certificate initially")
	}
	certificateProvider.certificateExpiry = certExpiry

	certChecker := newCertChecker(certificateProvider, make(chan bool, 1))
	go certChecker.refreshLoop()

	return certificateProvider, nil
}

func (cp *CertificateProvider) extractNewCertificate() error {
	certPEMBlock, err := os.ReadFile(cp.certFile)
	if err != nil {
		return errors.Wrap(err, "failed to read PEM block from certificate file")
	}
	keyPEMBlock, err := os.ReadFile(cp.keyFile)
	if err != nil {
		return errors.Wrap(err, "failed to read PEM block from key file")
	}
	// create a tls certificate to check, whether any error occurs
	_, err = tls.X509KeyPair(certPEMBlock, keyPEMBlock)
	if err != nil {
		return errors.Wrap(err, "failed to parse public/private key")
	}

	cp.mu.Lock()
	defer cp.mu.Unlock()
	cp.certPEMBlock = certPEMBlock
	cp.keyPEMBlock = keyPEMBlock
	return nil
}
