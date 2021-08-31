package main

import (
	"bioenpro4to_http_client/lib/bep4t_http_client"
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

func testBEP4TClient(statePsw string) (*string, error) {
	client := bep4t_http_client.NewBEP4THttpClient("192.168.1.91", 8000, false)
	cred, err := client.GetAuthCredential("aa000aa", "ciao", "did:iota:test:HRyXLr22JbT4VYczsFRB3p7T5xnHDReHU78d4Ns7RqAa")
	if err != nil {
		fmt.Printf("%s\n", err)
		return nil, err
	}
	fmt.Printf("%s\n", cred)

	err = client.IsCredentialValid(cred)
	if err == nil {
		fmt.Println("Credential is Valid")
	} else {
		fmt.Printf("%s\n", err)
	}

	err = client.NewDailyActorChannel(cred, statePsw, "01/09/2021")
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Println("Channel Created")
	}

	chanBase64, err := client.GetDailyActorChannel(cred, statePsw, "01/09/2021")
	if err != nil {
		fmt.Printf("%s", err)
		return nil, err
	}
	fmt.Printf("%s\n", *chanBase64)
	return chanBase64, nil
}

func testSendMsg() {
	statePsw := "This is my password"
	stateBase64, err := testBEP4TClient(statePsw)
	if err != nil {
		return
	}

	dailyCh, err := manager.NewDailyChannel(*stateBase64, statePsw)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer dailyCh.Drop()
	info := dailyCh.ChannelInfo()
	fmt.Printf("\n%s:%s\n\n\n", info.ChannelId, info.AnnounceId)

	keyNonce := manager.NewEncryptionKeyNonce("This is a secret key", "This is a secret nonce")
	public := newMessage("This is a public Message").toJson()
	private := newMessage("This is a private Message").toJson()

	packet := manager.NewRawPacket(public, private)
	msgId, err := dailyCh.SendRawPacket(packet, keyNonce)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Println("Message Sent: ", string(public))
	fmt.Println(*msgId)
}

func main() {

	config := idManager.NewPersistenceConfig("id_manager_backup", "psw")
	idManager, _ := idManager.NewIdentityManager(false, config)
	defer idManager.Drop()

	client := bep4t_http_client.NewBEP4THttpClient("192.168.1.91", 8000, false)
	//cred, _ := client.GetAuthCredential("aa000aa", "ciao", "did:iota:test:HRyXLr22JbT4VYczsFRB3p7T5xnHDReHU78d4Ns7RqAa")
	//
	//idManager.StoreCredential("cred1", cred)
	cred, _ := idManager.GetCredential("cred1")
	fmt.Println(string(cred))
	err := client.IsCredentialValid(cred)
	if err == nil {
		fmt.Println("Credential is valid")
	}
}
