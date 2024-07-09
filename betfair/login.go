package betfair

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"strings"
	"time"
)

// TODO test failing login, should i use omitempty on SessionToken
type loginResponse struct {
	LoginStatus  string `json:"loginStatus"`
	SessionToken string `json:"sessionToken"`
}

func (c *Client) Login() error {
	url := createUrl(login_url, "certlogin")

	response, err := doLoginRequest(c, url)
	if err != nil {
		return err
	}
	var result loginResponse
	err = json.Unmarshal(response, &result)
	if err != nil {
		return err
	}
	if result.LoginStatus != "SUCCESS" {
		return errors.New(fmt.Sprintf("login status not SUCCESS is  %s", result.LoginStatus))
	}

	// Update client with session information
	c.session.SessionToken = result.SessionToken
	c.session.LoginTime = time.Now().UTC()
	return nil
}

func doLoginRequest(c *Client, url string) ([]byte, error) {

	// Create tls.Config with the client certificate
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*c.certificates},
	}

	// Create http.Transport with the tls.Config
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Create http.Client with the transport
	client := &http.Client{
		Timeout:   5 * time.Second,
		Transport: transport,
	}

	// Create body
	body := neturl.Values{}
	body.Set("username", c.loginConfig.Username)
	body.Set("password", c.loginConfig.Password)

	request, err := http.NewRequest(http.MethodPost,
		url,
		strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Accept", "application/json")
	request.Header.Set("X-Application", c.loginConfig.AppKey)

	// Make a request
	resp, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid response code %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	// Handle response
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
