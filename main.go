package main

import (
	bep4t_client "bioenpro4to_http_client/lib"
	manager "bioenpro4to_http_client/lib/go_channel_manager"
	idManager "bioenpro4to_http_client/lib/identity_manager"
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
	config := idManager.NewPersistenceConfig("id_manager_backup", "psw")
	actorId := "aa000aa"
	//psw := "ciao"
	channelPsw := "psw"
	date := "02/09/2021"

	bep4tClient, err :=
		bep4t_client.BEP4TClientBuilder().HostAddr("192.168.1.91").Port(8000).MainNet(false).PersistenceConfig(config).Build()
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

	//err = bep4tClient.Authenticate(actorId, psw)
	//if err != nil{
	//	fmt.Println(err)
	//	return
	//}

	fmt.Println(bep4tClient.IsAuthenticated(actorId))
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
