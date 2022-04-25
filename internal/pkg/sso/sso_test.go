package sso

import (
	"log"
	"testing"
)

func TestClient_GenerateAccessToken(t *testing.T) {
	ssoClient := New("https://sso.bingosoft.net", "ajcNcUVYSmEW99qCyA9PnT", "b25da097-657d-4ed0-a579-47da34ad87e1", "bingo", "pass@cmp#2019")
	log.Println(ssoClient.GenerateAccessTokenByUser())
	log.Println(ssoClient.GenerateAccessTokenByClient())
}
