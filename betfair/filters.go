package betfair

import "time"

// important to use pointers here else omitempty will not work so
// to will always be included with the zero value and no response will match
type TimeRangeFilter struct {
	From *time.Time `json:"from,omitempty"`
	To   *time.Time `json:"to,omitempty"`
}

type MarketFilter struct {
	TextQuery          string          `json:"textQuery,omitempty"`
	EventTypeIds       []string        `json:"eventTypeIds,omitempty"`
	MarketCountries    []string        `json:"marketCountries,omitempty"`
	MarketIds          []string        `json:"marketIds,omitempty"`
	EventIds           []string        `json:"eventIds,omitempty"`
	CompetitionIds     []string        `json:"competitionIds,omitempty"`
	BSPOnly            bool            `json:"bspOnly,omitempty"`
	TurnInPlayEnabled  bool            `json:"turnInPLayEnabled,omitempty"`
	InPlayOnly         bool            `json:"inPlayOnly,omitempty"`
	MarketBettingTypes []string        `json:"marketBettingTypes,omitempty"`
	MarketTypeCodes    []string        `json:"marketTypeCodes,omitempty"`
	MarketStartTime    TimeRangeFilter `json:"marketStartTime,omitempty"`
	WithOrders         string          `json:"withOrders,omitempty"`
}

type Params struct {
	MarketFilter     MarketFilter `json:"filter,omitempty"`
	MaxResults       int          `json:"maxResults,omitempty"`
	Granularity      string       `json:"granularity,omitempty"`
	MarketProjection []string     `json:"marketProjection,omitempty"`
	Sort             string       `json:"sort,omitempty"`
	Locale           string       `json:"locale,omitempty"`
}
