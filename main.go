package main

import (
	"bioenpro4to_http_client/lib"
	"fmt"
)

func main(){
	client := lib.NewBEP4TClient("192.168.1.91", 8000, false)
	cred, err := client.GetAuthCredential("d111", "ciao", "did:iota:test:DLkyWU3jJFgK81KUB3YaDqkwQGMcFNYXTBzj8R4Qhopr")
	if err != nil{
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Printf("%s\n", cred)

	err = client.IsCredentialValid(cred)
	if err == nil{
		fmt.Println("Credential is Valid")
	} else{
		fmt.Printf("%s\n", err)
	}

	err = client.NewDailyActorChannel(cred, "psw", "01/08/2021")
	if err != nil{
		fmt.Printf("%s", err)
	}else{
		fmt.Println("Channel Created")
	}
}
