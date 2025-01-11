package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/anthdm/hollywood/actor"
	"github.com/gorilla/websocket"
	"github.com/khaleelsyed/game_server/types"
	"golang.org/x/exp/rand"
)

type PlayerSession struct {
	sessionID int
	clientID  int
	username  string
	conn      *websocket.Conn
}

type GameServer struct {
	ctx      *actor.Context
	sessions map[*actor.PID]struct{}
}

func newGameServer() actor.Receiver {
	return &GameServer{
		sessions: make(map[*actor.PID]struct{}),
	}
}

var upgrader websocket.Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func newPlayerSession(sid int, conn *websocket.Conn) actor.Producer {
	return func() actor.Receiver {
		return &PlayerSession{
			sessionID: sid,
			conn:      conn,
		}
	}
}

func (playerSession *PlayerSession) Receive(c *actor.Context) {
	switch c.Message().(type) {
	case actor.Started:
		playerSession.readLoop()
	}
}

func (playerSession *PlayerSession) readLoop() {
	var msg types.Message
	for {
		if err := playerSession.conn.ReadJSON(&msg); err != nil {
			fmt.Println("Read error: ", err)
			return
		}
		go playerSession.handleMessage(msg)
	}
}

func (playerSession *PlayerSession) handleMessage(msg types.Message) {
	switch msg.Type {
	case "login":
		var loginMsg types.LoginData
		if err := json.Unmarshal(msg.Data, &loginMsg); err != nil {
			panic(err)
		}
		playerSession.clientID = loginMsg.ClientID
		playerSession.username = loginMsg.Username
	case "playerState":
		var playerState types.PlayerState
		if err := json.Unmarshal(msg.Data, &playerState); err != nil {
			panic(err)
		}
		fmt.Println(playerState)
	}
}

func (s *GameServer) Receive(c *actor.Context) {
	switch msg := c.Message().(type) {
	case actor.Started:
		s.startHTTP()
		s.ctx = c
		_ = msg

	}
}

func (s *GameServer) startHTTP() {
	fmt.Println("Starting HTTP server on 40000")
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

	fmt.Println("New client attempting connection:")
	sid := rand.Intn(math.MaxInt)
	pid := s.ctx.SpawnChild(newPlayerSession(sid, conn), fmt.Sprintf("session_%d", sid))
	s.sessions[pid] = struct{}{}
	fmt.Printf("Client with sid %d and pid %s just connected\n", sid, pid)
}

func main() {
	e, err := actor.NewEngine(actor.NewEngineConfig())
	if err != nil {
		log.Fatal(err)
	}
	e.Spawn(newGameServer, "server")

	select {}
}

// Last at 1:27:00
