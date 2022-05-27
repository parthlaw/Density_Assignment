package websocket

import "encoding/json"

const PlaceBetAction = "place_bet"
const UpdateCounterAction = "update_counter"
const SyncAction = "sync"

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Sender  *Client `json:"sender"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		return nil
	}
	return json
}
