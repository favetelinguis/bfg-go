package main

import (
	"os"

	"github.com/favetelinguis/bfg-go/betfair"
)

func newLoginConfig() betfair.LoginConfig {
	appKey := os.Getenv("BFG_APP_KEY")
	password := os.Getenv("BFG_PASSWORD")
	username := os.Getenv("BFG_USERNAME")
	for _, str := range []string{appKey, password, username} {
		if str == "" {
			panic("Unset environment variable")
		}
	}

	return betfair.LoginConfig{
		AppKey:   appKey,
		Password: password,
		Username: username,
		CertFile: "betfair-2048.crt",
		KeyFile:  "betfair-2048.key",
	}
}
