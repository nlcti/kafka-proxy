package proxy

import (
	"crypto/tls"
	"time"

	"github.com/grepplabs/kafka-proxy/pkg/apis"
	"github.com/grepplabs/kafka-proxy/pkg/libs/util"
	"github.com/pkg/errors"
)

type TLSCertificateClient struct {
	certificate          tls.Certificate
	certificateAvailable bool
	certificateExpiry    time.Time
	certificateProvider  apis.CertificateProvider
}

func NewCertificateClient(certificateProvider apis.CertificateProvider) *TLSCertificateClient {
	return &TLSCertificateClient{certificateProvider: certificateProvider, certificateAvailable: false}
}

func (cc *TLSCertificateClient) GetCertificate(helloInfo *tls.ClientHelloInfo) (*tls.Certificate, error) {
	var reqExpiryDate string
	if cc.certificateExpiry.IsZero() {
		reqExpiryDate = time.Now().Format(util.CertTimeStampLayout)
	} else {
		reqExpiryDate = cc.certificateExpiry.Format(util.CertTimeStampLayout)
	}

	certPEMBlock, keyPEMBlock, err := cc.certificateProvider.GetX509KeyPair(reqExpiryDate)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request X509 key pair from certificate provider plugin")
	}
	// no new certificate available
	if certPEMBlock == nil || keyPEMBlock == nil {
		if cc.certificateAvailable {
			return &cc.certificate, nil
		} else {
			return nil, errors.New("no TLS certificate available")
		}
	} else {
		// received PEM data for new certificate
		cert, err := tls.X509KeyPair(certPEMBlock, keyPEMBlock)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse public/private key")
		}
		certExpiry, err := util.ExtractExpiryFromCertificate(certPEMBlock)
		if err != nil {
			return nil, errors.Wrap(err, "failed to extract certificate expiry")
		}
		cc.certificate = cert
		cc.certificateExpiry = certExpiry
		cc.certificateAvailable = true

		return &cc.certificate, nil
	}
}
