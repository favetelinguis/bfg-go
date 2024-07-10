package cache

type MarketChangeMessage struct {
	Op          string         `json:"op"`
	Id          int            `json:"id"`
	Ct          string         `json:"ct"`
	SegmentType string         `json:"segmentType"`
	ConflateMs  *int           `json:"conflateMs,omitempty"`
	Status      string         `json:"status"`
	HeartbeatMs int            `json:"heartbeatMs"`
	Pt          int64          `json:"pt"`
	InitialClk  string         `json:"initialClk"`
	Clk         string         `json:"clk"`
	Mc          []MarketChange `json:"mc"`
}

type MarketChange struct {
	Rc               []RunnerChange   `json:"rc"`
	Img              bool             `json:"img"`
	Tv               float32          `json:"tv"`
	MarketDefinition MarketDefinition `json:"marketDefinition"`
	Id               string           `json:"id"`
}

type RunnerChange struct {
	Id    int          `json:"id"`
	Con   bool         `json:"con"`
	Tv    float32      `json:"tv"`
	Ltp   float32      `json:"ltp"`
	Spn   float32      `json:"spn"`
	Spf   float32      `json:"spf"`
	Batb  [][3]float32 `json:"batb"`
	Batl  [][3]float32 `json:"batl"`
	Bdatb [][3]float32 `json:"bdatb"`
	Bdatl [][3]float32 `json:"bdatl"`
	Atb   [][2]float32 `json:"atb"`
	Atl   [][2]float32 `json:"atl"`
	Spb   [][2]float32 `json:"spb"`
	Spl   [][2]float32 `json:"spl"`
	Trd   [][2]float32 `json:"trd"`
	Hc    float32      `json:"hc"`
}

type MarketDefinition struct {
	Venue                 string     `json:"venue,omitempty"`
	BspMarket             bool       `json:"bspMarket,omitempty"`
	TurnInPlayEnabled     bool       `json:"turnInPlayEnabled,omitempty"`
	PersistenceEnabled    bool       `json:"persistenceEnabled,omitempty"`
	MarketBaseRate        float32    `json:"marketBaseRate,omitempty"`
	EventId               string     `json:"eventId,omitempty"`
	EventTypeId           string     `json:"eventTypeId,omitempty"`
	NumberOfWinners       int        `json:"numberOfWinners,omitempty"`
	BettingType           string     `json:"bettingType,omitempty"`
	MarketType            string     `json:"marketType,omitempty"`
	MarketTime            string     `json:"marketTime,omitempty"`
	SuspendTime           string     `json:"suspendTime,omitempty"`
	BspReconciled         bool       `json:"bspReconciled,omitempty"`
	Complete              bool       `json:"complete,omitempty"`
	InPlay                bool       `json:"inPlay,omitempty"`
	CrossMatching         bool       `json:"crossMatching,omitempty"`
	RunnersVoidable       bool       `json:"runnersVoidable,omitempty"`
	NumberOfActiveRunners int        `json:"numberOfActiveRunners,omitempty"`
	BetDelay              int        `json:"betDelay,omitempty"`
	Status                string     `json:"status,omitempty"`
	Runners               []McRunner `json:"runners,omitempty"`
	Regulators            []string   `json:"regulators,omitempty"`
	DiscountAllowed       bool       `json:"discountAllowed,omitempty"`
	Timezone              string     `json:"timezone,omitempty"`
	OpenDate              string     `json:"openDate,omitempty"`
	Version               int64      `json:"version,omitempty"`
}

type Runner struct {
	Status       string `json:"status,omitempty"`
	SortPriority int    `json:"sortPriority,omitempty"`
	Id           int    `json:"id,omitempty"`
}
