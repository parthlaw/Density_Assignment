package main

import (
	"density_assignment/binance"
	ws "density_assignment/websocket"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
)

func main() {
	// u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}
	// c, _, err := websocket.DefaultDialer.Dial("wss://stream.binance.com:9443/ws/ethusdt@kline_30m", nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer c.Close()
	// done := make(chan struct{})
	binan := binance.NewBinance()
	wsserver := ws.NewWebsocketServer()
	go func() {
		binan.ConnectToWebsocket(wsserver)
	}()
	go func() {
		ticker := time.NewTicker(60 * time.Second)
		for t := range ticker.C {
			binan.Timer <- t
		}
	}()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello World")) })
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) { fmt.Println("Request"); ws.ServeWs(wsserver, w, r) })
	c := cors.New(cors.Options{AllowedOrigins: []string{"http://localhost:3000", "https://bako.vercel.app", "https://gochat.vercel.app", "https://gochat.parthlaw.tech"}, AllowCredentials: true, AllowedHeaders: []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"}})
	handler := c.Handler(mux)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server Running on port %s\n", port)
	e := http.ListenAndServe(":"+port, handler)
	if e != nil {
		log.Fatal(e)
	}
}
