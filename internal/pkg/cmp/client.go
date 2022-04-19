package cmp

import (
	"github.com/levigross/grequests"
)

type config struct {
	Endpoint       string
	AccessToken    string
	MainApiContext string
	Options        grequests.RequestOptions
}

type Client struct {
	config *config
}

func New(endpoint, accessToken string) *Client {
	cmpClient := &Client{}
	cmpClient.config = &config{}
	cmpClient.config.Endpoint = endpoint
	cmpClient.config.AccessToken = accessToken
	cmpClient.config.MainApiContext = "gateway/cmp-main-api"

	var headers map[string]string
	if accessToken != "" {
		headers = map[string]string{"Authorization": "Bearer " + accessToken}
	}
	cmpClient.config.Options = grequests.RequestOptions{Headers: headers, InsecureSkipVerify: true}

	return cmpClient
}
