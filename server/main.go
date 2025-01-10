package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
)

type GameServer struct{}

type HTTPServer struct{}

func newGameServer() actor.Receiver {
	return &GameServer{}
}

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (s *GameServer) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		s.startHTTP()
		_ = msg

	}
}

func (s *GameServer) startHTTP() {
	fmt.Print("Starting HTTP server on 40000")
	go func() {
		http.HandleFunc("/ws", s.handleWS)
		http.ListenAndServe(":40000", nil)
	}()

}

// Handles the WS upgrade request
func (s *GameServer) handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("WS Upgrade error", err)
		return
	}

	fmt.Print("New client attempting connection:")
	fmt.Print(conn)
}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}
	e.Spawn(newGameServer, "server")

	select {}
}
