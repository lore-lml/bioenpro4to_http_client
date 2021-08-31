package bep4t_client

import (
	http "bioenpro4to_http_client/lib/bep4t_http_client"
	chManager "bioenpro4to_http_client/lib/go_channel_manager"
	idManager "bioenpro4to_http_client/lib/identity_manager"
)

type BEP4TClient struct {
	httpClient      *http.BEP4THttpClient
	identityManager *idManager.IdentityManager
	channels        map[string]*chManager.DailyChannel
	Authenticated   bool
}

func NewBEP4TClient(httpClient *http.BEP4THttpClient,
	identityManger *idManager.IdentityManager) *BEP4TClient {
	return &BEP4TClient{
		httpClient:      httpClient,
		identityManager: identityManger,
		channels:        make(map[string]*chManager.DailyChannel),
		Authenticated:   false,
	}
}
