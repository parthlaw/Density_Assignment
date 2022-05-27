package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Bet string

const (
	UpBet   Bet = "up"
	DownBet     = "down"
)

type Client struct {
	conn       *websocket.Conn
	wsServer   *WsServer
	send       chan []byte
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	Curr_bet   Bet       `json:"curr_bet"`
	Counter    int       `json:"counter"`
	Curr_Price string    `json:"curr_price"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

func newClient(conn *websocket.Conn, name string, wsServer *WsServer) *Client {
	return &Client{
		conn:     conn,
		wsServer: wsServer,
		send:     make(chan []byte, 256),
		ID:       uuid.New(),
		Name:     name,
		Curr_bet: UpBet,
		Counter:  0,
	}
}
func (client *Client) readPump() {
	defer func() {

	}()
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// client.wsServer.removeClient(client)
				client.conn.Close()
				log.Printf("Unexpected close error: %v", err)
			}
			break
		}
		client.handleNewMessages(jsonMessage)
	}
}
func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
func (client *Client) handleNewMessages(jsonMessage []byte) {
	var message Message
	err := json.Unmarshal(jsonMessage, &message)
	if err != nil {
		log.Printf("error: %v", err)
		return
	}
	switch message.Action {
	case PlaceBetAction:
		client.placeBet(message.Message)
	}
}
func (client *Client) placeBet(bet string) {
	if bet == "up" {
		client.Curr_bet = UpBet
	} else {
		client.Curr_bet = DownBet
	}
}
func (client *Client) UpdateCounter(price float32) {
	curr_price, _ := strconv.ParseFloat(client.Curr_Price, 32)
	if client.Curr_bet == UpBet {
		if price >= float32(curr_price) {
			client.Counter++
		} else {
			fmt.Println('\n', price, float32(curr_price), '\n')
			client.Counter--
		}
	} else {
		if price <= float32(curr_price) {
			client.Counter++
		} else {
			client.Counter--
		}
	}
}
func ServeWs(wsServer *WsServer, w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	if !ok || len(name[0]) < 1 {
		log.Println("Url params missing")
		return
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := newClient(conn, name[0], wsServer)
	fmt.Println("New Client Joined")
	fmt.Println(client)
	go client.writePump()
	go client.readPump()
	tm := int32(time.Now().Second()) - wsServer.start_time
	if tm < 0 {
		tm = 60 + tm
	}
	msg := &Message{
		Action:  SyncAction,
		Message: fmt.Sprint(tm),
	}
	wsServer.register <- client
	client.send <- msg.encode()
}
