package request_info

import (
	"bioenpro4to_http_client/lib/utils"
	"encoding/json"
	"errors"
)

type ChannelAuthorization struct {
	Cred       interface{} `json:"Cred"`
	ChannelPsw string           `json:"Channel-psw"`
}

func NewChannelAuthorization(cred utils.Credential, channelPsw string) (*ChannelAuthorization, error){
	var obj interface{}
	err := json.Unmarshal(cred, &obj)
	if err != nil{
		return nil, errors.New("Invalid credential format")
	}
	return &ChannelAuthorization{
		Cred: obj,
		ChannelPsw: channelPsw,
	}, nil
}

