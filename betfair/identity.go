package betfair

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type identityResponse struct {
	Token   string `json:"token"`
	Product string `json:"product"`
	Status  string `json:"status"`
	Error   string `json:"error"`
}

func (s *Session) Logout() error {
	if s.token.SessionToken == "" {
		return nil
	}
	url := createUrl(identity_url, "logout")

	response, err := doIdentityRequest(s, url)
	if err != nil {
		return err
	}

	var result identityResponse
	if err = json.Unmarshal(response, &result); err != nil {
		return err
	}

	if result.Status != "SUCCESS" {
		return fmt.Errorf("logout status not SUCCESS is  %s with reason %s", result.Status, result.Error)
	}

	s.token.SessionToken = ""
	s.token.LoginTime = time.Time{}
	return nil
}

func (c *Session) KeepAlive() error {
	url := createUrl(identity_url, "keepAlive")

	response, err := doIdentityRequest(c, url)
	if err != nil {
		return err
	}

	var result identityResponse
	if err = json.Unmarshal(response, &result); err != nil {
		return err
	}

	if result.Status != "SUCCESS" {
		return fmt.Errorf("keepAlive status not SUCCESS is  %s with reason %s", result.Status, result.Error)
	}

	c.token.SessionToken = result.Token
	c.token.LoginTime = time.Now().UTC()

	return nil
}

func doIdentityRequest(c *Session, url string) ([]byte, error) {

	r, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("X-Application", c.loginConfig.AppKey)
	r.Header.Set("X-Authentication", c.token.SessionToken)

	client := &http.Client{}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid response code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
