package request_info

import (
	"bioenpro4to_http_client/lib/utils"
)

type ChannelAuthorization struct {
	Cred       string `json:"Cred"`
	ChannelPsw string `json:"Channel-psw"`
}

func NewChannelAuthorization(cred utils.Credential, channelPsw string) *ChannelAuthorization {
	return &ChannelAuthorization{
		Cred:       string(cred),
		ChannelPsw: channelPsw,
	}
}

func (self *ChannelAuthorization) ToMap() map[string]string {
	m := make(map[string]string)
	m["Cred"] = self.Cred
	m["Channel-psw"] = self.ChannelPsw
	return m
}
