package betfair

import (
	"context"
	"fmt"
	"time"

	"github.com/ybbus/jsonrpc/v3"
)

type EventType struct {
	Id   *string `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
}

type Competition struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type Event struct {
	Id          string `json:"id,omitempty"`
	OpenDate    string `json:"openDate,omitempty"`
	TimeZone    string `json:"timezone,omitempty"`
	CountryCode string `json:"countryCode,omitempty"`
	Name        string `json:"name,omitempty"`
	Venue       string `json:"venue,omitempty"`
}

type MarketCatalogueDescription struct {
	BettingType        string    `json:"bettingType"`
	BSPMarket          bool      `json:"bspMarket"`
	DiscountAllowed    bool      `json:"discountAllowed"`
	MarketBaseRate     float32   `json:"marketBaseRate"`
	MarketTime         time.Time `json:"marketTime"`
	MarketType         string    `json:"marketType"`
	PersistenceEnabled bool      `json:"persistenceEnabled"`
	Regulator          string    `json:"regulator"`
	SuspendDate        time.Time `json:"suspendTime"`
	TurnInPlayEnabled  bool      `json:"turnInPlayEnabled"`
	Rules              *string   `json:"rules,omitempty"`
	RulesHasDate       *bool     `json:"rulesHasDate,omitempty"`
	Wallet             *string   `json:"wallet,omitempty"`
	EachWayDivisor     *float32  `json:"eachWayDivisor,omitempty"`
	Clarifications     *string   `json:"clarifications,omitempty"`
}

type Metadata struct {
	RunnerId int `json:"runnerId"`
}

type RunnerCatalogue struct {
	SelectionId  int     `json:"selectionId"`
	RunnerName   string  `json:"runnerName"`
	SortPriority int     `json:"sortPriority"`
	Handicap     float32 `json:"handicap"`
	//Metadata		*metadata	`json:"metadata"`  //TODO
}

type MarketCatalogue struct {
	MarketId                   string                     `json:"marketId"`
	MarketName                 string                     `json:"marketName"`
	TotalMatched               *float32                   `json:"totalMatched,omitempty"`
	MarketStartTime            *time.Time                 `json:"marketStartTime,omitempty"`
	Competition                Competition                `json:"competition,omitempty"`
	Event                      Event                      `json:"event,omitempty"`
	EventType                  EventType                  `json:"eventType,omitempty"`
	MarketCatalogueDescription MarketCatalogueDescription `json:"description,omitempty"`
	Runners                    []RunnerCatalogue          `json:"runners,omitempty"`
}

func (b *Betting) ListMarketCatalogue() ([]MarketCatalogue, error) {
	from := time.Now().UTC()
	params := Params{}
	params.MarketFilter = MarketFilter{
		EventTypeIds:    []string{"7"},
		MarketCountries: []string{"GB"},
		MarketTypeCodes: []string{"WIN"},
		MarketStartTime: TimeRangeFilter{
			From: &from,
		},
	}
	params.Sort = "FIRST_TO_START"
	params.MaxResults = 10
	params.MarketProjection = []string{"MARKET_START_TIME"}

	methodName := "listMarketCatalogue"
	client := jsonrpc.NewClientWithOpts(api_betting_url, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"X-Application":    b.Session.loginConfig.AppKey,
			"X-Authentication": b.Session.token.SessionToken,
			// TODO should i use keep-alive?
			// "Connection": "keep-alive",
		},
		AllowUnknownFields: false,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var response []MarketCatalogue
	err := client.CallFor(
		ctx,
		&response,
		createJsonRpcMethodName(methodName),
		params,
	)
	if err != nil {
		return nil, fmt.Errorf("jsonrpc call failed: %w", err)
	}

	return response, nil
}
