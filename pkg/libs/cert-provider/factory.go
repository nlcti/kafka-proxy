package certprovider

import (
	"flag"

	"github.com/grepplabs/kafka-proxy/pkg/apis"
	"github.com/grepplabs/kafka-proxy/pkg/registry"
)

func init() {
	registry.NewComponentInterface(new(apis.CertificateProviderFactory))
	registry.Register(new(Factory), "cert-provider")
}

type pluginMeta struct {
	certFile         string
	keyFile          string
	checkIntervalMin int
}

func (f *pluginMeta) flagSet() *flag.FlagSet {
	fs := flag.NewFlagSet("certificate provider plugin settings", flag.ContinueOnError)

	fs.StringVar(&f.certFile, "updated-proxy-listener-cert-file", "", "PEM encoded file with server certificate, being updated frequently")
	fs.StringVar(&f.keyFile, "updated-proxy-listener-key-file", "", "PEM encoded file with private key for the server certificate")
	fs.IntVar(&f.checkIntervalMin, "update-check-interval-minutes", 0, "Interval in minutes to check for an update")

	return fs
}

type Factory struct {
}

func (f *Factory) New(params []string) (apis.CertificateProvider, error) {

	pluginMeta := &pluginMeta{}
	flags := pluginMeta.flagSet()
	flags.Parse(params)

	options := CertificateProviderOptions{
		CertFile:         pluginMeta.certFile,
		KeyFile:          pluginMeta.keyFile,
		CheckIntervalMin: pluginMeta.checkIntervalMin,
	}

	return NewCertificateProvider(options)
}
