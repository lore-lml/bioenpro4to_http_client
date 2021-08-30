package request_info

import "strings"

type AuthInfo struct {
	Id string `json:"id"`
	Psw string `json:"psw"`
	Did string `json:"did"`
}

func NewAuthInfo(id, psw, did string) *AuthInfo{
	return &AuthInfo{
		Id: strings.ToLower(id),
		Psw: psw,
		Did: did,
	}
}
