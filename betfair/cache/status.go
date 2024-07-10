package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

type StatusMessage struct {
	Op                   string  `json:"op"`
	Id                   int     `json:"id"`
	StatusCode           string  `json:"statusCode"`
	ConnectionClosed     bool    `json:"connectionClosed"`
	ErrorCode            *string `json:"errorCode,omitempty"`
	ErrorMessage         *string `json:"errorMessage,omitempty"`
	ConnectionsAvailable *int    `json:"connectionsAvailable,omitempty"`
}

type StatusCache struct {
	IsConnectionClosed   bool
	ConnectionsAvailable int
	ResponseChans        map[int]chan bool
	Mu                   sync.Mutex
}

func NewStatusCache() *StatusCache {
	return &StatusCache{
		ResponseChans: make(map[int]chan bool),
	}
}

func (s *StatusCache) Parse(message string) error {
	var statusMessage StatusMessage

	err := json.Unmarshal([]byte(message), &statusMessage)
	if err != nil {
		return fmt.Errorf("unable to parse: %w", err)
	}

	s.IsConnectionClosed = statusMessage.ConnectionClosed
	if statusMessage.ConnectionsAvailable != nil {
		s.ConnectionsAvailable = *statusMessage.ConnectionsAvailable
	}

	s.Mu.Lock()
	defer s.Mu.Unlock()

	if ch, ok := s.ResponseChans[statusMessage.Id]; ok {
		var Err error
		switch statusMessage.StatusCode {
		case "SUCCESS":
			ch <- true
		case "FAILURE":
			ch <- false
			if statusMessage.ErrorMessage != nil {
				Err = errors.New(*statusMessage.ErrorMessage)
			} else {
				Err = errors.New("betfair response status is FAILURE")
			}
		default:
			ch <- false
			Err = fmt.Errorf("unknown statuscode %s", statusMessage.StatusCode)
		}
		close(ch)
		delete(s.ResponseChans, statusMessage.Id)
		return Err
	}
	return nil
}
