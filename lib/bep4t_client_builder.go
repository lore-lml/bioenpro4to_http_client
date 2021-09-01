package bep4t_client

import (
	"bioenpro4to_http_client/lib/bep4t_http_client"
	"bioenpro4to_http_client/lib/identity_manager"
)

type bep4tClientBuilder struct {
	hostAddr          string
	port              int16
	mainNet           bool
	persistenceConfig *identity_manager.PersistenceConfig
}

func BEP4TClientBuilder() *bep4tClientBuilder {
	return &bep4tClientBuilder{
		hostAddr:          "localhost",
		port:              8080,
		persistenceConfig: nil,
	}
}

func (self *bep4tClientBuilder) HostAddr(hostAddr string) *bep4tClientBuilder {
	self.hostAddr = hostAddr
	return self
}

func (self *bep4tClientBuilder) Port(port int16) *bep4tClientBuilder {
	self.port = port
	return self
}

func (self *bep4tClientBuilder) MainNet(mainnet bool) *bep4tClientBuilder {
	self.mainNet = mainnet
	return self
}

func (self *bep4tClientBuilder) PersistenceConfig(config *identity_manager.PersistenceConfig) *bep4tClientBuilder {
	self.persistenceConfig = config
	return self
}

func (self *bep4tClientBuilder) Build() (*BEP4TClient, error) {
	httpClient := bep4t_http_client.NewBEP4THttpClient(self.hostAddr, self.port, false)
	identityManager, err := identity_manager.NewIdentityManager(self.mainNet, self.persistenceConfig)
	if err != nil {
		return nil, err
	}

	return newBEP4TClient(httpClient, identityManager), nil
}
