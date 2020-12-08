package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/fabiofcferreira/novaengelapiparser"
	"github.com/fatih/color"
)

// Login performs login request
func Login(username string, password string) (*novaengelapiparser.LoginAuthorization, error) {
	authorization := &novaengelapiparser.LoginAuthorization{}

	requestBody, err := json.Marshal(map[string]string{
		"user":     username,
		"password": password,
	})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(links["login"], "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		color.HiRed("Couldn't perform login request.")
		return nil, err
	}
	defer resp.Body.Close()

	// Parse request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		color.HiRed("Couldn't parse response JSON.")
		return nil, err
	}

	// Unmarshal JSON
	err = json.Unmarshal(body, authorization)
	if err != nil {
		color.HiRed("Couldn't parse response JSON.")
		return nil, err
	}

	return authorization, err
}
