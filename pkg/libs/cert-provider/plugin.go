package certprovider

import (
	"crypto/tls"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/grepplabs/kafka-proxy/pkg/libs/util"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type CertificateProvider struct {
	mu                sync.RWMutex
	certFile          string
	keyFile           string
	certPEMBlock      []byte
	keyPEMBlock       []byte
	certificateExpiry time.Time
}

type CertificateProviderOptions struct {
	CertFile string
	KeyFile  string
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

	certificateProvider := &CertificateProvider{certFile: opts.CertFile, keyFile: opts.KeyFile}

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

	handleFileUpdates := func() {
		newCertExpiry, err := util.ExtractExpiryFromCertificateFile(certificateProvider.certFile)
		if err != nil {
			logrus.Error(fmt.Sprintf("failed to extract certificate expiry - %v", err))
			return
		}
		if !newCertExpiry.IsZero() && newCertExpiry.After(certificateProvider.certificateExpiry) {
			err = certificateProvider.extractNewCertificate()
			if err != nil {
				logrus.Error(fmt.Sprintf("failed to create certificate - %v", err))
				return
			}
			certificateProvider.certificateExpiry = newCertExpiry
		}
	}

	stopChannelCert := make(chan bool, 1)
	err = util.WatchForUpdates(certificateProvider.certFile, stopChannelCert, handleFileUpdates)
	if err != nil {
		return nil, errors.Wrap(err, "cannot watch certificate file")
	}
	stopChannelKey := make(chan bool, 1)
	err = util.WatchForUpdates(certificateProvider.keyFile, stopChannelKey, handleFileUpdates)
	if err != nil {
		return nil, errors.Wrap(err, "cannot watch key file")
	}

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
