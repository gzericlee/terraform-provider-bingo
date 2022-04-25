package sso

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/levigross/grequests"

	"terraform-provider-bingo/utils"
)

type Authorization struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func (its Authorization) String() string {
	return utils.Prettify(its)
}

func (its *Client) GenerateAccessTokenByClient() (*Authorization, error) {
	options := its.config.Options
	resp, err := grequests.Post(fmt.Sprintf("%v/oauth2/token?grant_type=client_credentials", its.config.Endpoint), &options)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	content := resp.String()
	if !resp.Ok {
		err = fmt.Errorf("[SSO] Response code: [%v]，result: [%s]", resp.StatusCode, content)
		return nil, err
	}

	auth := &Authorization{}
	err = json.Unmarshal([]byte(content), &auth)

	return auth, err
}

func (its *Client) GenerateAccessTokenByUser() (*Authorization, error) {
	options := its.config.Options
	resp, err := grequests.Post(fmt.Sprintf("%v/oauth2/token?grant_type=password&username=%v&password=%v", its.config.Endpoint, its.config.UserName, url.QueryEscape(its.config.Password)), &options)
	if err != nil {
		return nil, err
	}
	defer resp.Close()

	content := resp.String()
	if !resp.Ok {
		err = fmt.Errorf("[SSO] Response code: [%v]，result: [%s]", resp.StatusCode, content)
		return nil, err
	}

	auth := &Authorization{}
	err = json.Unmarshal([]byte(content), &auth)

	return auth, err
}
