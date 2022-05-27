package websocket

import (
	"fmt"
	"time"
)

type WsServer struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	start_time int32
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		start_time: int32(time.Now().Second()),
	}
}
func (ws *WsServer) registerClient(client *Client) {
	ws.clients[client] = true
}
func (ws *WsServer) unregisterClient(client *Client) {
	if _, ok := ws.clients[client]; ok {
		delete(ws.clients, client)
		close(client.send)
	}
}
func (ws *WsServer) broadcastToClients(message []byte) {
	for client := range ws.clients {
		client.send <- message
	}
}
func (ws *WsServer) Run() {
	for {
		select {
		case client := <-ws.register:
			ws.registerClient(client)
			// log.Printf("Client %s registered", client.ID)
		case client := <-ws.unregister:
			ws.unregisterClient(client)
			// log.Printf("Client %s unregistered", client.ID)
		case message := <-ws.broadcast:
			ws.broadcastToClients(message)
			// log.Printf("Message sent: %s", message)
		}
	}
}
func (ws *WsServer) UpdateAllCounter(price float32) {
	for client := range ws.clients {
		client.UpdateCounter(price)
		mstr := fmt.Sprintf("Curr Price is %f, Counter is %d", price, client.Counter)
		message := &Message{
			Action:  UpdateCounterAction,
			Message: mstr,
		}
		client.send <- message.encode()
	}
}
