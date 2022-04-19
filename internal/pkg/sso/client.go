package sso

import "github.com/levigross/grequests"

type config struct {
	Endpoint     string
	ClientSecret string
	Options      grequests.RequestOptions
}

type Client struct {
	config *config
}

func New(endpoint, clientSecret string) *Client {
	ssoClient := &Client{}
	ssoClient.config = &config{}
	ssoClient.config.Endpoint = endpoint
	ssoClient.config.ClientSecret = clientSecret

	var headers map[string]string
	if clientSecret != "" {
		headers = map[string]string{"Authorization": "Basic " + clientSecret}
	}
	ssoClient.config.Options = grequests.RequestOptions{Headers: headers, InsecureSkipVerify: true}

	return ssoClient
}
