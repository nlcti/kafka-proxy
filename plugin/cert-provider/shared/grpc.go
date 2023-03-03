package shared

import (
	"github.com/grepplabs/kafka-proxy/pkg/apis"
	"github.com/grepplabs/kafka-proxy/plugin/cert-provider/proto"
	"github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of CertificateProvider that talks over gRPC.
type GRPCClient struct {
	broker *plugin.GRPCBroker
	client proto.CertificateProviderClient
}

func (m *GRPCClient) GetX509KeyPair(expiryDate string) ([]byte, []byte, error) {
	resp, err := m.client.GetX509KeyPair(context.Background(), &proto.CertificateRequest{ExpiryDate: expiryDate})
	if err != nil {
		return nil, nil, err
	}

	return resp.CertPEMBlock, resp.KeyPEMBlock, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	proto.UnimplementedCertificateProviderServer
	broker *plugin.GRPCBroker
	Impl   apis.CertificateProvider
}

func (m *GRPCServer) GetX509KeyPair(ctx context.Context, req *proto.CertificateRequest) (*proto.CertificateResponse, error) {
	certPEMBlock, keyPEMBlock, err := m.Impl.GetX509KeyPair(req.ExpiryDate)

	return &proto.CertificateResponse{CertPEMBlock: certPEMBlock, KeyPEMBlock: keyPEMBlock}, err
}
