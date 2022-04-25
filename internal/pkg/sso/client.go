package sso

import (
	"encoding/base64"
	"fmt"

	"github.com/levigross/grequests"
)

type config struct {
	Endpoint     string
	ClientId     string
	ClientSecret string
	Options      grequests.RequestOptions
}

type Client struct {
	config *config
}

func New(endpoint, clientId, clientSecret string) *Client {
	ssoClient := &Client{}
	ssoClient.config = &config{}
	ssoClient.config.Endpoint = endpoint
	ssoClient.config.ClientId = clientId
	ssoClient.config.ClientSecret = clientSecret

	var headers map[string]string
	if clientSecret != "" {
		headers = map[string]string{
			"Authorization": "Basic " + base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%v:%v", clientId, clientSecret))),
		}
	}
	ssoClient.config.Options = grequests.RequestOptions{Headers: headers, InsecureSkipVerify: true}

	return ssoClient
}
