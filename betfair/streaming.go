package betfair

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/favetelinguis/bfg-go/betfair/cache"
)

type message interface {
	SetId(int)
}

type authenticationMessage struct {
	Op      string `json:"op"`
	Id      int    `json:"id"`
	AppKey  string `json:"appKey"`
	Session string `json:"session"`
}

type marketFilter struct {
	CountryCodes      []string `json:"countryCodes,omitempty"`
	BettingTypes      []string `json:"bettingTypes,omitempty"`
	TurnInPlayEnabled *bool    `json:"turnInPlayEnabled,omitempty"`
	MarketTypes       []string `json:"marketTypes,omitempty"`
	Venues            []string `json:"venues,omitempty"`
	MarketIds         []string `json:"marketIds,omitempty"`
	EventTypeIds      []string `json:"eventTypeIds,omitempty"`
	EventIds          []string `json:"eventIds,omitempty"`
	BspMarket         *bool    `json:"bspMarket,omitempty"`
	RaceTypes         []string `json:"raceTypes,omitempty"`
}

type marketDataFilter struct {
	LadderLevels *int     `json:"ladderLevels,omitempty"`
	Fields       []string `json:"fields,omitempty"`
}

type marketSubscriptionMessage struct {
	Op                  string           `json:"op"`
	Id                  int              `json:"id"`
	SegmentationEnabled *bool            `json:"segmentationEnabled,omitempty"`
	ConflateMs          *int             `json:"conflateMs,omitempty"`
	HeartbeatMs         *int             `json:"heartbeatMs,omitempty"`
	InitialClk          *string          `json:"initialClk,omitempty"`
	Clk                 *string          `json:"clk,omitempty"`
	MarketFilter        marketFilter     `json:"marketFilter"`
	MarketDataFilter    marketDataFilter `json:"marketDataFilter"`
}

func (msg *marketSubscriptionMessage) SetId(id int) {
	msg.Id = id
}

func (msg *authenticationMessage) SetId(id int) {
	msg.Id = id
}

func (s *Streaming) Connect() error {
	conf := &tls.Config{
		InsecureSkipVerify: false,
	}

	conn, err := tls.Dial("tcp", stream_url, conf)
	if err != nil {
		return fmt.Errorf("stream connection failure: %w", err)
	}

	s.conn = conn
	s.closeCh = make(chan struct{})
	s.StatusCache = cache.NewStatusCache()
	go s.receiveLoop()
	return nil
}

func (s *Streaming) Close() {
	// TODO not sure how well this is working
	// looks like I close the connection and get an error in the
	// receive loop while i would expect channel would close first to exit
	// receive loop before conn i closed?
	if s.conn != nil {
		close(s.closeCh)
		// time.Sleep(1 * time.Second) // do not help
		s.conn.Close()
		s.conn = nil
	}
}

func (s *Streaming) Authenticate() error {
	if s.conn == nil {
		return fmt.Errorf("Unable to authenticate since the stream is not connected")
	}
	authMessage := &authenticationMessage{
		Op:      "authentication",
		AppKey:  s.Session.loginConfig.AppKey,
		Session: s.Session.token.SessionToken,
	}

	err := s.send(authMessage)
	if err != nil {
		return fmt.Errorf("failed to send auth message: %w", err)
	}
	return nil
}

func (s *Streaming) SubscribeToMarkets(marketIds []string) error {
	mf := &marketFilter{
		MarketIds: marketIds,
	}

	mdf := &marketDataFilter{
		// TODO would prob want EX_BEST_OFFERS here to not get full ladder only top10
		// however i want the volume for all levels
		Fields: []string{"EX_ALL_OFFERS", "EX_TRADED", "EX_TRADED_VOL", "EX_LTP", "EX_MARKET_DEF", "EX_BEST_OFFERS_DISP"},
	}

	msm := &marketSubscriptionMessage{
		Op:               "marketSubscription",
		MarketFilter:     *mf,
		MarketDataFilter: *mdf,
	}

	err := s.send(msm)
	if err != nil {
		return fmt.Errorf("failure seding market subscription: %w", err)
	}

	return nil
}

func (s *Streaming) parse(message string) error {
	var msgMap map[string]interface{}

	err := json.Unmarshal([]byte(message), &msgMap)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	op, exists := msgMap["op"]
	if !exists {
		return errors.New("missing 'op' field in message")
	}
	opStr, ok := op.(string)
	if !ok {
		return errors.New("'op' field is not a string")
	}

	switch opStr {
	case "status":
		err := s.StatusCache.Parse(message)
		if err != nil {
			return fmt.Errorf("unable to update status cache: %w", err)
		}
	case "connection":
		fmt.Println(message)
	case "mcm":
		fmt.Println(message)
	case "ocm":
		fmt.Println(message)
	default:
		return fmt.Errorf("onknown 'op': %s", opStr)
	}
	return nil
}

func (s *Streaming) nextMsgId() int {
	s.msgCount++
	return s.msgCount
}

func (s *Streaming) send(msg message) error {
	if s.conn == nil {
		return fmt.Errorf("unable to send, not connected")
	}

	id := s.nextMsgId()
	msg.SetId(id)

	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal send data: %w", err)
	}
	bytes = append(bytes, "\r\n"...)

	responseChan := make(chan bool, 1)
	s.StatusCache.Mu.Lock()
	s.StatusCache.ResponseChans[id] = responseChan
	s.StatusCache.Mu.Unlock()

	_, err = s.conn.Write(bytes)
	if err != nil {
		s.StatusCache.Mu.Lock()
		delete(s.StatusCache.ResponseChans, id)
		s.StatusCache.Mu.Unlock()
		return fmt.Errorf("sending failed %w", err)
	}

	// block until we get a response or fail after 2 seconds
	select {
	case status := <-responseChan:
		if status {
			return nil
		} else {
			return fmt.Errorf("failure when sending")
		}
	case <-time.After(2 * time.Second):
		s.StatusCache.Mu.Lock()
		delete(s.StatusCache.ResponseChans, id)
		s.StatusCache.Mu.Unlock()
		return fmt.Errorf("sending timed out")
	}
}

func (s *Streaming) receiveLoop() {
	reader := bufio.NewReader(s.conn)
	for {
		select { // TODO need to understand the select, is not im blocking on readstring so the closeCh will never work?
		case <-s.closeCh:
			return
		default:
			response, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("Error receiving data: %v", err)
				return
			}
			err = s.parse(response)
			if err != nil {
				log.Printf("unable to parse message: %s", response)
			}
		}
	}
}
