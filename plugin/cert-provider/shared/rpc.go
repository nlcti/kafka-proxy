package shared

import (
	"net/rpc"

	"github.com/grepplabs/kafka-proxy/pkg/apis"
)

type RPCClient struct {
	client *rpc.Client
}

func (m *RPCClient) GetX509KeyPair(expiryDate string) ([]byte, []byte, error) {
	var resp map[string]interface{}

	err := m.client.Call("Plugin.GetX509KeyPair", map[string]interface{}{
		"expiryDate": expiryDate,
	}, &resp)

	return resp["certPEMBlock"].([]byte), resp["keyPEMBlock"].([]byte), err
}

type RPCServer struct {
	Impl apis.CertificateProvider
}

func (m *RPCServer) GetX509KeyPair(args map[string]interface{}, resp *map[string]interface{}) error {
	certPEMBlock, keyPEMBlock, err := m.Impl.GetX509KeyPair(args["expiryDate"].(string))

	*resp = map[string]interface{}{
		"certPEMBlock": certPEMBlock,
		"keyPEMBlock":  keyPEMBlock,
	}

	return err
}
