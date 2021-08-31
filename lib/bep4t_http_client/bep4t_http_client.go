package bep4t_http_client

import (
	"bioenpro4to_http_client/lib/bep4t_http_client/request_info"
	"bioenpro4to_http_client/lib/bep4t_http_client/utils"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	authenticateURL    = "id-manager/authenticate"
	validCredentialURL = "id-manager/is-credential-valid"
	dailyChannelURL    = "channel-manager/daily-channel"
)

type Method string

const (
	httpGet  Method = "GET"
	httpPost        = "POST"
)

type BEP4THttpClient struct {
	BaseUrl    string
	httpClient *http.Client
}

func NewBEP4THttpClient(hostAddr string, port int16, ssl bool) *BEP4THttpClient {
	var secure = ""
	if ssl {
		secure = "s"
	}

	return &BEP4THttpClient{
		BaseUrl:    fmt.Sprintf("http%s://%s:%s", secure, hostAddr, strconv.FormatInt(int64(port), 10)),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (client *BEP4THttpClient) newRequest(method Method, apiUrl string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", client.BaseUrl, apiUrl)
	req, err := http.NewRequest(string(method), url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (client *BEP4THttpClient) Welcome() {
	req, err := client.newRequest(httpGet, "", nil)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	res, err := client.httpClient.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return
	}
	if res.StatusCode > 299 {
		fmt.Printf(string(body))
	}
	fmt.Printf("%s", body)
}

func (client *BEP4THttpClient) GetAuthCredential(actorId, psw, did string) (utils.Credential, error) {
	auth, err := json.Marshal(request_info.NewAuthInfo(actorId, psw, did))
	if err != nil {
		return nil, err
	}

	req, err := client.newRequest(httpPost, authenticateURL, bytes.NewBuffer(auth))
	if err != nil {
		return nil, err
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	cred, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		return nil, errors.New(string(cred))
	}
	return cred, nil
}

func (client *BEP4THttpClient) IsCredentialValid(cred utils.Credential) error {
	req, err := client.newRequest(httpGet, validCredentialURL, bytes.NewBuffer(cred))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode > 299 {
		return errors.New(string(body))
	}

	success := struct {
		IsValid bool `json:"is_valid"`
	}{IsValid: false}
	err = json.Unmarshal(body, &success)
	if err != nil {
		return err
	}

	if !success.IsValid {
		return errors.New("invalid credential")
	}

	return nil
}

func (client *BEP4THttpClient) NewDailyActorChannel(cred utils.Credential, channelPsw, date string) error {
	timestamp, err := utils.DateToTimestamp(date)
	if err != nil {
		return err
	}

	auth := request_info.NewChannelAuthorization(cred, channelPsw)

	strTimestamp := strconv.FormatInt(timestamp, 10)
	reqBody := []byte(fmt.Sprintf(`{"day_timestamp": %s}`, strTimestamp))
	req, err := client.newRequest(httpPost, dailyChannelURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return err
	}

	for key, value := range auth.ToMap() {
		req.Header.Add(key, value)
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if res.StatusCode > 299 {
		return errors.New(string(body))
	}
	return nil
}

func (client *BEP4THttpClient) GetDailyActorChannel(cred utils.Credential, channelPsw, date string) (*string, error) {
	if !utils.CheckDateFormat(date) {
		return nil, errors.New("wrong date format")
	}
	date = strings.ReplaceAll(date, "/", "-")
	auth := request_info.NewChannelAuthorization(cred, channelPsw)

	req, err := client.newRequest(httpGet, fmt.Sprintf("%s/%s", dailyChannelURL, date), nil)
	if err != nil {
		return nil, err
	}

	for key, value := range auth.ToMap() {
		req.Header.Add(key, value)
	}

	res, err := client.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode > 299 {
		return nil, errors.New(string(body))
	}

	js := struct {
		ChannelBase64 string `json:"channel_base64"`
	}{ChannelBase64: ""}
	err = json.Unmarshal(body, &js)
	if err != nil {
		return nil, err
	}

	return &js.ChannelBase64, nil
}
