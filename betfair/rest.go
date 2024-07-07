package betfair

import "time"

type LogoutResponse struct {
	Token   string
	Product string
	Status  string
	Error   string
}

type MarketCatalogue struct {
	MarketId        string
	MarketName      string
	MarketStartTime time.Time
}
