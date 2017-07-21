package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

type GetApiResponse struct {
	Count    int    `json:"count"`
	Next     string `next`
	Previous string `previous`
	List     []Api  `list`
}

type Api struct {
	Id          string `json:id`
	Name        string `json:name`
	Description string `json:description`
	Context     string `json:context`
	Version     string `json:version`
	Provider    string `json:provider`
}

func GetClientIdSecret(clientRegEndpoint string, username string, password string) (clientID string, clientSecret string) {

	client := createHTTPClient()
	var payload = []byte(`{
        "callbackUrl": "https://localhost/callback",
        "clientName": "wso2-apim-cli",
        "tokenScope": "Production",
        "owner": "admin",
        "grantType": "password refresh_token",
        "saasApp": true
        }`)

	req, err := http.NewRequest("POST", clientRegEndpoint, bytes.NewBuffer(payload))
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return "", err.Error()
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}

	clientID = data["clientId"].(string)
	clientSecret = data["clientSecret"].(string)

	log.Print("Client ID: ", clientID)
	log.Print("Client Secret: ", clientSecret)
	return clientID, clientSecret
}

func GetToken(tokenEndpoint string, username string, password string, clientId string, clientSecret string) (accessToken string) {
	client := createHTTPClient()
	req, err := http.NewRequest("POST", tokenEndpoint, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(clientId, clientSecret))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	q := req.URL.Query()
	q.Add("grant_type", "password")
	q.Add("username", username)
	q.Add("password", password)
	q.Add("scope", "apim:api_view apim:api_create apim:api_publish")
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var data map[string]interface{}
	err = json.Unmarshal([]byte(body), &data)
	if err != nil {
		panic(err)
	}

	accessToken = data["access_token"].(string)
	log.Println("Access token: ", accessToken)
	return accessToken
}

func ExportAPI(exportEndpoint string, username string, password string, exportPath string, name string, version string, provider string) (err error) {

	client := createHTTPClient()
	req, err := http.NewRequest("GET", exportEndpoint, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))

	q := req.URL.Query()
	q.Add("name", name)
	q.Add("version", version)
	q.Add("provider", provider)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode == http.StatusFound {
		return errors.New("API import/export web application not found")
	}
	if res.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(body))
	}

	out, err := os.Create(exportPath + "/" + name + "-" + version + ".zip")
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, res.Body)
	return nil
}

func ImportAPI(importEndpoint string, username string, password string, filePath string) (err error) {

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	formFile, err := writer.CreateFormFile("file", filePath)
	if err != nil {
		return
	}
	if _, err = io.Copy(formFile, file); err != nil {
		return
	}
	writer.Close()

	req, err := http.NewRequest("POST", importEndpoint, &buffer)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := createHTTPClient()
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == http.StatusFound {
		return errors.New("API import/export web application not found")
	}
	if res.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(body))
	}
	return
}

func GetAPIs(publisherEndpoint string, token string) GetApiResponse {

	client := createHTTPClient()
	req, err := http.NewRequest("GET", publisherEndpoint+"/v0.11/apis/", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	res, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return GetApiResponse{}
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var response GetApiResponse
	json.Unmarshal(body, &response)
	return response
}

var RedirectAttemptedError = errors.New("redirect")

func createHTTPClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, CheckRedirect: func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}}
	return client
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
