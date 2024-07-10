package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type McRunner struct {
	Id              int
	FullPriceLadder map[string]map[float32]float32
	SingleValues    map[string]float32
}

func NewRunner(id int) *McRunner {
	return &McRunner{
		Id:              id,
		FullPriceLadder: make(map[string]map[float32]float32),
		SingleValues:    make(map[string]float32),
	}
}

func (r *McRunner) Update(runnerChange RunnerChange) {

	r.UpdateFullPriceLadder(runnerChange.Atb, "atb")
	r.UpdateFullPriceLadder(runnerChange.Atl, "atl")
	r.UpdateFullPriceLadder(runnerChange.Trd, "trd")
	r.UpdateFullPriceLadder(runnerChange.Spb, "spb")
	r.UpdateFullPriceLadder(runnerChange.Spl, "spl")

	r.UpdateSingleValue(runnerChange.Tv, "tv")
	r.UpdateSingleValue(runnerChange.Ltp, "ltp")
	r.UpdateSingleValue(runnerChange.Spn, "spn")
	r.UpdateSingleValue(runnerChange.Spf, "spf")
	r.UpdateSingleValue(runnerChange.Hc, "hc")
}

func (r *McRunner) UpdateFullPriceLadder(list [][2]float32, selection string) {
	if list == nil {
		return
	}

	if _, ok := r.FullPriceLadder[selection]; !ok {
		r.FullPriceLadder[selection] = make(map[float32]float32)
	}

	for _, i := range list {
		price := i[0]
		size := i[1]
		r.FullPriceLadder[selection][price] = size
	}
}

func (r *McRunner) UpdateLevelBasedLadder(list [][3]float32, selection string) {
	if list == nil {
		return
	}

}

func (r *McRunner) UpdateSingleValue(singleValue float32, selection string) {
	if singleValue != 0 {
		r.SingleValues[selection] = singleValue
	}
}

type Market struct {
	Id               string
	Runners          map[int]*McRunner
	TotalVolume      int
	MarketDefinition MarketDefinition
}

func NewMarket(id string) *Market {
	return &Market{
		Id:      id,
		Runners: make(map[int]*McRunner),
	}
}

func (m *Market) Update(mc MarketChange) {
	if mc.Rc != nil {
		for _, runnerChange := range mc.Rc {
			id := runnerChange.Id
			_, ok := m.Runners[id]
			if !ok {
				m.AddRunner(id)
			}
			m.Runners[id].Update(runnerChange)

		}
	}
}

func (m *Market) AddRunner(id int) {
	runner := NewRunner(id)
	m.Runners[id] = runner
}

type MarketCache struct {
	Markets            map[string]*Market
	HeartbeatThreshold time.Duration
	InitialClk         string
	Clk                string
	mu                 sync.Mutex
	timer              *time.Timer
}

func NewMarketCache() *MarketCache {
	return &MarketCache{
		Markets:            make(map[string]*Market),
		HeartbeatThreshold: time.Duration(5000) * time.Millisecond,
	}
}

func (m *MarketCache) Parse(msg string) ([]Market, error) {
	var mcm MarketChangeMessage

	err := json.Unmarshal([]byte(msg), &mcm)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	timeSent := time.Unix(0, mcm.Pt*int64(time.Millisecond))
	timeSinceSent := time.Since(timeSent)

	if mcm.Status != "" {
		return nil, fmt.Errorf(mcm.Status)
	}

	if mcm.Clk != "" {
		m.Clk = mcm.Clk
	}

	if mcm.InitialClk != "" {
		m.InitialClk = mcm.InitialClk
	}

	if mcm.HeartbeatMs != 0 {
		m.HeartbeatThreshold = time.Duration(int64(float64(mcm.HeartbeatMs)*1.05)) * time.Millisecond
	}

	switch mcm.Ct {
	case "SUB_IMAGE":
		fmt.Print("")
	case "RESUB_DELTA":
		fmt.Print("")
	case "HEARTBEAT":
		fmt.Print("")
	case "":
		fmt.Print("")
	}

	var updatedMarkets []Market
	if mcm.Mc != nil {
		for _, marketChange := range mcm.Mc {
			id := marketChange.Id
			market, ok := m.Markets[id]
			if ok {
				market.Update(marketChange)
			} else {
				// This is the first message for the market so add to market cache
				// TODO this is very bad prob, better way?
				m.AddMarket(id)
				market, ok = m.Markets[id]
				if ok {
					market.Update(marketChange)
				}
			}
			updatedMarkets = append(updatedMarkets, *market)
		}
	}

	if timeSinceSent > 0 {
		return nil, fmt.Errorf("high latency %v", timeSinceSent)
	}

	m.resetTimer()

	return updatedMarkets, nil
}

func (m *MarketCache) resetTimer() {
	if m.timer != nil {
		m.timer.Stop()
	}
	m.timer = time.AfterFunc(m.HeartbeatThreshold, func() {
		m.mu.Lock()
		defer m.mu.Unlock()
		fmt.Println("Error: heartbeat missed")
	})
}

func (m *MarketCache) StopTimer() {
	if m.timer != nil {
		m.timer.Stop()
	}
}

func (m *MarketCache) AddMarket(id string) {
	market := NewMarket(id)
	m.Markets[id] = market
}
