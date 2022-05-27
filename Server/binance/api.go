package binance

import (
	wbs "density_assignment/websocket"
	"encoding/json"
	"fmt"
	"strconv"

	// "fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Binance struct {
	Timer chan time.Time
}

func NewBinance() *Binance {
	return &Binance{
		Timer: make(chan time.Time),
	}
}
func connectToBinance() *websocket.Conn {
	// Connect to Binance
	ws, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/ws/ethusdt@kline_1m", nil)
	if err != nil {
		log.Fatal("Error", err)
	}
	return ws
}
func (binance *Binance) ConnectToWebsocket(wsServer *wbs.WsServer) {
	// Connect to Binance
	ws := connectToBinance()
	defer ws.Close()
	// Start the server
	go wsServer.Run()
	// Listen for messages from the Binance API
	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Println("read err", err)
			return
		}
		// log.Printf("recv: %s", message)
		t := <-binance.Timer
		fmt.Print(t.Second())
		// log.Printf("recv: %s", message)
		handleNewMessages(message, wsServer)
	}
}
func handleNewMessages(jsonMessage []byte, wsserver *wbs.WsServer) {
	var message Message
	err := json.Unmarshal(jsonMessage, &message)
	if err != nil {
		log.Printf("error: %v", err)
	}
	fmt.Print(message)
	op, _ := strconv.ParseFloat(message.Kline.Open_Price, 32)
	wsserver.UpdateAllCounter(float32(op))
}
