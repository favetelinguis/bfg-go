package betfair

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ybbus/jsonrpc/v3"
)

type Session struct {
	token            string
	client           http.Client
	rpcClient        jsonrpc.RPCClient
	connectionConfig connectionConfig
}

// TODO schedule a keep-alive rutine that starts running once my session is created
func NewSession() *Session {
	conf := newConnectinConfig()
	sessionToken := login(conf.appKey, conf.username, conf.password)

	return &Session{
		token: sessionToken,
		client: http.Client{
			Timeout: 5 * time.Second,
		},
		// TODO test if betfair.se work
		rpcClient: jsonrpc.NewClientWithOpts("https://api.betfair.com/exchange/betting/json-rpc/v1", &jsonrpc.RPCClientOpts{
			CustomHeaders: map[string]string{
				"X-Application":    conf.appKey,
				"X-Authentication": sessionToken,
			},
			AllowUnknownFields: true,
		}),
		connectionConfig: conf,
	}
}

// Handle all the retry and error checking for requests,
// TODO think this should only be used for keep-alive and logout?
func (p *Session) connectedClient(url string, method string, body io.Reader) ([]byte, error) {
	r, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("Accept", "application/json")
	r.Header.Set("X-Application", p.connectionConfig.appKey)
	r.Header.Set("X-Authentication", p.token)

	response, err := p.client.Do(r)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode > 299 {
		return nil, fmt.Errorf("Invalid response code %d", response.StatusCode)
	}

	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return responseBody, nil
}

func (p *Session) Logout() *LogoutResponse {
	body, err := p.connectedClient("https://identitysso.betfair.se/api/logout", http.MethodPost, nil)
	if err != nil {
		panic(err)
	}

	var res LogoutResponse
	if err = json.Unmarshal(body, &res); err != nil {
		panic(err)
	}
	fmt.Printf("Logout with: %v", res)
	return &res
}

func (p *Session) KeepAlive() (*LogoutResponse, error) {
	body, err := p.connectedClient("https://identitysso.betfair.se/api/keepAlive", http.MethodPost, nil)
	if err != nil {
		return nil, err
	}

	var res LogoutResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (p *Session) CallListMarketCatalogue() ([]MarketCatalogue, error) {
	var data []MarketCatalogue
	err := p.rpcClient.CallFor(
		context.Background(), // TODO how to use context here
		&data,
		"SportsAPING/v1.0/listMarketCatalogue",
		map[string]interface{}{
			"filter": map[string]interface{}{
				"eventTypeIds":    []string{"7"},
				"marketCountries": []string{"GB"},
				"marketTypeCodes": []string{"WIN"},
				"marketStartTime": map[string]interface{}{
					"from": time.Now().UTC()},
			},
			"sort":             "FIRST_TO_START",
			"maxResults":       "10",
			"marketProjection": []string{"MARKET_START_TIME"},
		},
	)

	if err != nil || data == nil {
		return nil, err
	}
	return data, nil
}
