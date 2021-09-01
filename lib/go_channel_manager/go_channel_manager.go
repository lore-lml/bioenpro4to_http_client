package go_channel_manager

/*
#cgo LDFLAGS: -L./.. -lc_channel_manager_lib
#include "../c_channel_manager.h"
*/
import "C"
import (
	"errors"
	"unsafe"
)

func Hello(str string) {
	cString := C.CString(str)
	C.hello_from_rust(cString)
}

type DailyChannel struct {
	channel *C.daily_channel_t
}

type ChannelInfo struct {
	ChannelId  string
	AnnounceId string
}

type KeyNonce struct {
	key   string
	nonce string
}

type RawPacket struct {
	Public []byte
	Masked []byte
}

func NewDailyChannel(stateBase64, statePsw string) (*DailyChannel, error) {
	cStateBase64 := C.CString(stateBase64)
	cStatePsw := C.CString(statePsw)
	dailyChannel := C.daily_channel_from_base64(cStateBase64, cStatePsw)
	if dailyChannel == nil {
		return nil, errors.New("bad state format or wrong password")
	}
	return &DailyChannel{channel: dailyChannel}, nil
}

func (ch *DailyChannel) Drop() {
	C.drop_daily_channel_manager(ch.channel)
}

func (ch *DailyChannel) SendRawPacket(packet *RawPacket, keyNonce *KeyNonce) (string, error) {
	var kn *C.key_nonce_t = nil
	if keyNonce != nil {
		kn = keyNonce.toCKeyNonce()
		defer C.drop_key_nonce(kn)
	}

	pack := packet.toCRawPacket()
	defer C.drop_raw_packet(pack)

	cMsgId := C.send_raw_packet(ch.channel, pack, kn)
	if cMsgId == nil {
		return "", errors.New("something wrong during sending the packet")
	}
	defer C.drop_str(cMsgId)
	msgId := C.GoString(cMsgId)
	return msgId, nil
}

func (ch *DailyChannel) ChannelInfo() *ChannelInfo {
	info := C.daily_channel_info(ch.channel)
	defer C.drop_channel_info(info)
	return NewChannelInfo(C.GoString(info.channel_id), C.GoString(info.announce_id))
}

func NewChannelInfo(channelId, announceId string) *ChannelInfo {
	return &ChannelInfo{ChannelId: channelId, AnnounceId: announceId}
}

func NewEncryptionKeyNonce(key, nonce string) *KeyNonce {
	return &KeyNonce{key: key, nonce: nonce}
}

func (keyNonce *KeyNonce) toCKeyNonce() *C.key_nonce_t {
	return C.new_encryption_key_nonce(C.CString(keyNonce.key), C.CString(keyNonce.nonce))
}

func NewRawPacket(pubData, maskData []byte) *RawPacket {
	return &RawPacket{Public: pubData, Masked: maskData}
}

func (packet *RawPacket) toCRawPacket() *C.raw_packet_t {
	cPub, pLen := goByteToCByte(packet.Public)
	cMask, mLen := goByteToCByte(packet.Masked)
	return C.new_raw_packet(cPub, pLen, cMask, mLen)
}

func goByteToCByte(bytes []byte) (*C.uchar, C.ulong) {
	if bytes == nil {
		return nil, 0
	}
	return (*C.uchar)(unsafe.Pointer(&bytes[0])), C.ulong(len(bytes))
}
