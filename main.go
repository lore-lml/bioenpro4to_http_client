package main

import (
	"bioenpro4to_http_client/lib"
	"fmt"
)

func main() {
	client := lib.NewBEP4TClient("192.168.1.91", 8000, false)
	cred, err := client.GetAuthCredential("aa000aa", "ciao", "did:iota:test:HRyXLr22JbT4VYczsFRB3p7T5xnHDReHU78d4Ns7RqAa")
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("%s\n", cred)

	err = client.IsCredentialValid(cred)
	if err == nil {
		fmt.Println("Credential is Valid")
	} else {
		fmt.Printf("%s\n", err)
	}

	err = client.NewDailyActorChannel(cred, "psw", "31/08/2021")
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Println("Channel Created")
	}

	chanBase64, err := client.GetDailyActorChannel(cred, "psw", "31/08/2021")
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	fmt.Printf("%s", *chanBase64)
}
