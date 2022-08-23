package models

import (
	"encoding/json"
)

type WebSocketInput struct {
	Method string                 `json:"method"`
	Value  map[string]interface{} `json:"value"`
}

type WebSocketResponse struct {
	Data map[string]interface{} `json:"data"`
}

func (wsr *WebSocketResponse) ToResponse() []byte {
	data, err := json.Marshal(wsr)
	if err != nil {
		return []byte("no response")
	}
	return data
}

type CompetitorVote struct {
	Competitor int `json:"competitor"`
}
