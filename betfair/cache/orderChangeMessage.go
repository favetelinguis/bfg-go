package cache

type OrderChangeMessage struct {
	Op          string               `json:"op"`
	Id          int                  `json:"id"`
	Ct          string               `json:"ct"`
	SegmentType string               `json:"segmentType"`
	ConflateMs  *int                 `json:"conflateMs,omitempty"`
	Status      string               `json:"status"`
	HeartbeatMs int                  `json:"heartbeatMs"`
	Pt          int64                `json:"pt"`
	InitialClk  string               `json:"initialClk"`
	Clk         string               `json:"clk"`
	Oc          []OrderAccountChange `json:"oc"`
}

type OrderAccountChange struct {
	Closed    bool          `json:"closed"`
	ID        string        `json:"id"`
	FullImage bool          `json:"fullImage"`
	Orc       []OrderChange `json:"orc"`
}

type OrderChange struct {
	FullImage bool             `json:"fullImage"`
	Id        string           `json:"id"`
	Hc        *string          `json:"hc,omitempty"`
	Uo        []UnmatchedOrder `json:"uo"`
	Mb        [][]float32      `json:"mb"`
	Ml        [][]float32      `json:"ml"`
}

type UnmatchedOrder struct {
	Id     string  `json:"id"`
	P      float32 `json:"p"`
	S      float32 `json:"s"`
	Bsp    float32 `json:"bsp"`
	Side   string  `json:"side"`
	Status string  `json:"status"`
	Pt     string  `json:"pt"`
	Ot     string  `json:"ot"`
	Pd     int     `json:"pd"`
	Md     int     `json:"md"`
	Cd     int     `json:"cd"`
	Ld     int     `json:"ld"`
	Lsrc   string  `json:"lsrc"`
	Avp    float32 `json:"avp"`
	Sm     float32 `json:"sm"`
	Sr     float32 `json:"sr"`
	Sl     float32 `json:"sl"`
	Sc     float32 `json:"sc"`
	Sv     float32 `json:"sv"`
	Rac    string  `json:"rac"`
	Rc     string  `json:"rc"`
	Rfo    string  `json:"rfo"`
	Rfs    string  `json:"rfs"`
}
