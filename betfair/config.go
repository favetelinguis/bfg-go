package betfair

import "os"

type connectionConfig struct {
	appKey   string
	password string
	url      string
	username string
}

func newConnectinConfig() connectionConfig {
	appKey := os.Getenv("BFG_APP_KEY")
	password := os.Getenv("BFG_PASSWORD")
	url := os.Getenv("BFG_URL")
	username := os.Getenv("BFG_USERNAME")
	for _, str := range []string{appKey, password, url, username} {
		if str == "" {
			panic("Unset environment variable")
		}
	}

	return connectionConfig{
		appKey:   appKey,
		password: password,
		url:      url,
		username: username,
	}
}
