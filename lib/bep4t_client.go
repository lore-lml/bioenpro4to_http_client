package bep4t_client

import (
	http "bioenpro4to_http_client/lib/bep4t_http_client"
	chManager "bioenpro4to_http_client/lib/go_channel_manager"
	idManager "bioenpro4to_http_client/lib/identity_manager"
	"errors"
	"fmt"
	"strings"
)

type BEP4TClient struct {
	httpClient      *http.BEP4THttpClient
	identityManager *idManager.IdentityManager
	channels        map[string]*chManager.DailyChannel
}

func newBEP4TClient(httpClient *http.BEP4THttpClient, identityManger *idManager.IdentityManager) *BEP4TClient {
	return &BEP4TClient{
		httpClient:      httpClient,
		identityManager: identityManger,
		channels:        make(map[string]*chManager.DailyChannel),
	}
}

func (self *BEP4TClient) CreateIdentity(actorId string) (string, error) {
	return self.identityManager.CreateIdentity(actorId)
}

func (self *BEP4TClient) GetIdentityDid(actorId string) (string, error) {
	return self.identityManager.GetIdentityDid(actorId)
}

func (self *BEP4TClient) Authenticate(actorId, bep4tPsw string) error {
	if self.IsAuthenticated(actorId) {
		return nil
	}

	did, err := self.GetIdentityDid(actorId)
	if err != nil {
		return err
	}

	cred, err := self.httpClient.GetAuthCredential(actorId, bep4tPsw, did)
	if err != nil {
		return err
	}

	err = self.identityManager.StoreCredential(actorId, "ch-auth", cred)
	if err != nil {
		return err
	}

	return nil
}

func (self *BEP4TClient) IsAuthenticated(actorId string) bool {
	cred, err := self.identityManager.GetCredential(actorId, "ch-auth")
	if err != nil {
		return false
	}
	err = self.httpClient.IsCredentialValid(cred)
	return err == nil
}

func (self *BEP4TClient) NewDailyChannel(actorId, channelPsw, date string) error {
	cred, err := self.identityManager.GetCredential(actorId, "ch-auth")
	if err != nil {
		return err
	}
	err = self.httpClient.NewDailyActorChannel(cred, channelPsw, date)
	if err != nil {
		return err
	}

	return self.RestoreDailyChannel(actorId, channelPsw, date)
}

func (self *BEP4TClient) RestoreDailyChannel(actorId, channelPsw, date string) error {
	cred, err := self.identityManager.GetCredential(actorId, "ch-auth")
	if err != nil {
		return err
	}

	chBase64, err := self.httpClient.GetDailyActorChannel(cred, channelPsw, date)
	if err != nil {
		return err
	}

	dailyCh, err := chManager.NewDailyChannel(chBase64, channelPsw)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("%s:%s", actorId, strings.ReplaceAll(date, "-", "/"))
	self.channels[key] = dailyCh

	return nil
}

func (self *BEP4TClient) SendRawPacketToChannel(actorId, date string, packet *chManager.RawPacket, keyNonce *chManager.KeyNonce) (string, error) {
	key := fmt.Sprintf("%s:%s", actorId, strings.ReplaceAll(date, "-", "/"))
	dailyCh, ok := self.channels[key]
	if !ok {
		return "", errors.New(fmt.Sprintf("daily channel in date %s not found, create it or try to restore it", date))
	}

	return dailyCh.SendRawPacket(packet, keyNonce)
}

func (self *BEP4TClient) InfoOfChannel(actorId, date string) (*chManager.ChannelInfo, error) {
	key := fmt.Sprintf("%s:%s", actorId, strings.ReplaceAll(date, "-", "/"))
	dailyCh, ok := self.channels[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("daily channel in date %s not found, create it or try to restore it", date))
	}
	return dailyCh.ChannelInfo(), nil
}

func (self *BEP4TClient) Drop() {
	for _, ch := range self.channels {
		ch.Drop()
	}
}
