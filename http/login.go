package http

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/fabiofcferreira/novaengelapiscraper"
	"github.com/fatih/color"
)

// Login performs login request
func Login(username string, password string) (*novaengelapiscraper.LoginAuthorization, error) {
	authorization := &novaengelapiscraper.LoginAuthorization{}

	requestBody, err := json.Marshal(map[string]string{
		"user":     username,
		"password": password,
	})
	if err != nil {
		log.Fatalln(err)
	}

	resp, err := http.Post(Links["login"], "application/json", bytes.NewBuffer(requestBody))
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

	// Add login time
	authorization.LastLoggedIn = time.Now()

	return authorization, err
}
