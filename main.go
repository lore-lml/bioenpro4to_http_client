package main

import (
	bep4t_client "bioenpro4to_http_client/lib"
	"bioenpro4to_http_client/lib/env_configuration"
	manager "bioenpro4to_http_client/lib/go_channel_manager"
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func newMessage(message string) *Message {
	return &Message{Message: message, Timestamp: time.Now().Unix()}
}

func (self *Message) toJson() []byte {
	res, _ := json.Marshal(self)
	return res
}

func main() {
	env, err := env_configuration.InitEnvConfiguration()
	if err != nil {
		fmt.Println(err)
		return
	}

	config := env.IdentityConfig
	actorId := env.ActorId
	actorAuthPsw := env.ActorAuthPsw
	channelPsw := env.ActorChannelPsw
	date := "03/09/2021"

	bep4tClient, err := bep4t_client.BEP4TClientBuilder().HostAddr(env.HostAddr).Port(env.HostPort).MainNet(env.Mainnet).PersistenceConfig(config).Build()
	defer bep4tClient.Drop()

	if err != nil {
		fmt.Println(err)
		return
	}

	_, err = bep4tClient.GetIdentityDid(actorId)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bep4tClient.Authenticate(actorId, actorAuthPsw)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(bep4tClient.IsAuthenticated(actorId))

	err = bep4tClient.NewDailyChannel(actorId, channelPsw, date)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bep4tClient.RestoreDailyChannel(actorId, channelPsw, date)
	if err != nil {
		fmt.Println(err)
		return
	}

	info, _ := bep4tClient.InfoOfChannel(actorId, date)
	fmt.Printf("%s:%s\n", info.ChannelId, info.AnnounceId)

	//keyNonce := manager.NewEncryptionKeyNonce("This is a secret key", "This is a secret nonce")
	public := newMessage("This is a public Message").toJson()
	//private := newMessage("This is a private Message").toJson()
	packet := manager.NewRawPacket(public, nil)

	msgId, err := bep4tClient.SendRawPacketToChannel(actorId, date, packet, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(msgId)
}
