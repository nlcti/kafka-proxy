package certprovider

import (
	"fmt"
	"time"

	"github.com/grepplabs/kafka-proxy/pkg/libs/util"
	"github.com/sirupsen/logrus"
)

type certChecker struct {
	certProvider *CertificateProvider
	stopChannel  chan bool
	interval     time.Duration
}

func newCertChecker(certProvider *CertificateProvider, stopChannel chan bool) *certChecker {
	return &certChecker{
		certProvider: certProvider,
		stopChannel:  stopChannel,
		interval:     time.Duration(certProvider.checkIntervalMin) * time.Minute,
	}
}

func (cr *certChecker) refreshLoop() {
	defer func() {
		if r := recover(); r != nil {
			var ok bool
			err, ok := r.(error)
			if ok {
				logrus.Errorf("certificate refresh loop error %v", err)
			}
		}
	}()
	logrus.Infof("Checking for new certificate every: %v", cr.interval)
	syncTicker := time.NewTicker(cr.interval)
	for {
		select {
		case <-syncTicker.C:
			cr.checkForRefresh()
		case <-cr.stopChannel:
			return
		}
	}
}

func (cr *certChecker) checkForRefresh() {
	newCertExpiry, err := util.ExtractExpiryFromCertificateFile(cr.certProvider.certFile)
	if err != nil {
		logrus.Error(fmt.Sprintf("failed to extract certificate expiry - %v", err))
		return
	}
	if !newCertExpiry.IsZero() && newCertExpiry.After(cr.certProvider.certificateExpiry) {
		err = cr.certProvider.extractNewCertificate()
		if err != nil {
			logrus.Error(fmt.Sprintf("failed to create certificate - %v", err))
			return
		}
		cr.certProvider.certificateExpiry = newCertExpiry
	}
}
