package apis

type CertificateProvider interface {
	GetX509KeyPair(expiryDate string) ([]byte, []byte, error)
}

type CertificateProviderFactory interface {
	New(params []string) (CertificateProvider, error)
}
