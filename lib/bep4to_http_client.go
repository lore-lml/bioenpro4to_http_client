package lib

import (
	"bioenpro4to_http_client/lib/request_info"
	"bioenpro4to_http_client/lib/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)



type BEP4TClient struct{
	BaseUrl string
	httpClient *http.Client
}

func NewBEP4TClient(hostAddr string, port int16, ssl bool) *BEP4TClient {
	var secure = ""
	if ssl{
		secure = "s"
	}

	return &BEP4TClient{
		BaseUrl: fmt.Sprintf("http%s://%s:%s", secure, hostAddr, strconv.FormatInt(int64(port), 10)),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (client *BEP4TClient) newRequest(method utils.Method, apiUrl string, body io.Reader) (*http.Request, error){
	url := fmt.Sprintf("%s/%s", client.BaseUrl, apiUrl)
	req, err := http.NewRequest(string(method), url, body)
	if err != nil{
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (client *BEP4TClient) Welcome(){
	req, err := client.newRequest(utils.Get, "", nil)
	if err != nil{
		fmt.Printf("%s", err)
		return
	}
	res, err := client.httpClient.Do(req)
	if err != nil{
		fmt.Printf("%s", err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		fmt.Printf("%s", err)
		return
	}
	if res.StatusCode > 299{
		fmt.Printf(string(body))
	}
	fmt.Printf("%s", body)
}

func (client *BEP4TClient) GetAuthCredential(actorId, psw, did string) (utils.Credential, error){
	auth, err := json.Marshal(request_info.NewAuthInfo(actorId, psw, did))
	if err != nil{
		return nil, err
	}

	req, err := client.newRequest(utils.Post, utils.Authenticate, bytes.NewBuffer(auth))
	if err != nil{
		return nil, err
	}

	res, err := client.httpClient.Do(req)
	if err != nil{
		return nil, err
	}

	defer res.Body.Close()
	cred, err := ioutil.ReadAll(res.Body)
	if err != nil{
		return nil, err
	}
	if res.StatusCode > 299{
		return nil, errors.New(string(cred))
	}
	return cred, nil
}

func (client *BEP4TClient) IsCredentialValid(cred utils.Credential) error{
	req, err := client.newRequest(utils.Get, utils.ValidCredential, bytes.NewBuffer(cred))
	if err != nil{
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.httpClient.Do(req)
	if err != nil{
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		return err
	}
	if res.StatusCode > 299{
		return errors.New(string(body))
	}

	var objMap map[string]json.RawMessage
	err = json.Unmarshal(body, &objMap)
	if err != nil{
		return err
	}

	var success bool
	err = json.Unmarshal(objMap["is_valid"], &success)

	return err
}

func (client *BEP4TClient) NewDailyActorChannel(cred utils.Credential, channelPsw, date string) error{
	timestamp, err := utils.DateToTimestamp(date)
	if err != nil{
		return err
	}

	auth, err := request_info.NewChannelAuthorization(cred, channelPsw)
	if err != nil{
		return err
	}
	serAuth, err := json.Marshal(auth)
	if err != nil{
		return err
	}

	strTimestamp := strconv.FormatInt(timestamp, 10)
	reqBody := []byte(fmt.Sprintf(`{"day_timestamp": %s}`, strTimestamp))
	req, err := client.newRequest(utils.Post, utils.DailyChannel, bytes.NewBuffer(reqBody))
	if err != nil{
		return err
	}

	var objMap map[string]json.RawMessage
	err = json.Unmarshal(serAuth, &objMap)
	if err != nil{
		return err
	}

	for key, value := range  objMap{
		req.Header.Add(key, fmt.Sprintf("%s", value))
	}

	res, err := client.httpClient.Do(req)
	if err != nil{
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil{
		return err
	}
	if res.StatusCode > 299{
		return errors.New(string(body))
	}
	return nil
}
