package main

import (
	"log"
	"math"

	"github.com/gorilla/websocket"
	"golang.org/x/exp/rand"
)

const wsServerEndpoint = "ws://localhost:40000"

type LoginData struct {
	ClientID int    `json:"client_id"`
	Username string `json:"username"`
}

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	Username string
}

func (c *GameClient) login() error {
	return c.conn.WriteJSON(LoginData{
		ClientID: c.clientID,
		Username: c.Username,
	})
}

func newGameClient(conn *websocket.Conn, username string) *GameClient {
	return &GameClient{
		conn:     conn,
		clientID: rand.Intn(math.MaxInt),
		Username: username,
	}
}

func main() {
	dialer := websocket.Dialer{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, _, err := dialer.Dial(wsServerEndpoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	c := newGameClient(conn, "James")
	if err := c.login(); err != nil {
		log.Fatal(err)
	}
}
