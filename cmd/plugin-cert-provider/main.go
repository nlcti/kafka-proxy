package main

import (
	"os"

	certprovider "github.com/grepplabs/kafka-proxy/pkg/libs/cert-provider"
	"github.com/grepplabs/kafka-proxy/plugin/cert-provider/shared"
	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
)

func main() {
	certProvider, err := new(certprovider.Factory).New(os.Args[1:])
	if err != nil {
		logrus.Errorf("cannot initialize certificate provider: %v", err)
		os.Exit(1)
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"certificateProvider": &shared.CertificateProviderPlugin{Impl: certProvider},
		},
		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
