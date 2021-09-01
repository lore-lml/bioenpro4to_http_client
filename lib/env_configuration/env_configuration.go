package env_configuration

import (
	"bioenpro4to_http_client/lib/identity_manager"
	"errors"
	"github.com/joho/godotenv"
	"strconv"
)

type EnvConfiguration struct {
	ActorId         string
	ActorAuthPsw    string
	ActorChannelPsw string
	HostAddr        string
	HostPort        int16
	Mainnet         bool
	IdentityConfig  *identity_manager.PersistenceConfig
}

func InitEnvConfiguration() (*EnvConfiguration, error) {
	myEnv, err := godotenv.Read()
	if err != nil {
		return nil, errors.New("error loading .env file")
	}

	actorId, ok := myEnv["ACTOR.ID"]
	if !ok {
		return nil, errors.New("Missing ACTOR.ID variable in the .env file")
	}

	actorAuthPsw, ok := myEnv["ACTOR.AUTH_PSW"]
	if !ok {
		return nil, errors.New("Missing ACTOR.AUTH_PSW variable in the .env file")
	}

	actorChannelPsw, ok := myEnv["ACTOR.CHANNEL_PSW"]
	if !ok {
		return nil, errors.New("Missing ACTOR.CHANNEL_PSW variable in the .env file")
	}

	hostAddr, ok := myEnv["HOST.ADDR"]
	if !ok {
		hostAddr = "localhost"
	}

	portStr, ok := myEnv["HOST.PORT"]
	if !ok {
		portStr = "8000"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8000
	}

	var mainNet bool
	if myEnv["IOTA.MAINNET"] == "true" {
		mainNet = true
	} else {
		mainNet = false
	}

	var persistence *identity_manager.PersistenceConfig = nil
	if myEnv["IOTA.IDENTITY.STORAGE.TYPE"] == "stronghold" {
		folder, ok1 := myEnv["IOTA.IDENTITY.STORAGE.DIR"]
		psw, ok2 := myEnv["IOTA.IDENTITY.STORAGE.PSW"]
		if ok1 && ok2 {
			persistence = identity_manager.NewPersistenceConfig(folder, psw)
		}
	}

	return &EnvConfiguration{
		ActorId:         actorId,
		ActorAuthPsw:    actorAuthPsw,
		ActorChannelPsw: actorChannelPsw,
		HostAddr:        hostAddr,
		HostPort:        int16(port),
		Mainnet:         mainNet,
		IdentityConfig:  persistence,
	}, nil
}
