package main

import (
	"encoding/json"
	"log"
	"math"
	"time"

	"github.com/gorilla/websocket"
	"github.com/khaleelsyed/game_server/types"
	"golang.org/x/exp/rand"
)

const wsServerEndpoint = "ws://localhost:40000/ws"

type GameClient struct {
	conn     *websocket.Conn
	clientID int
	Username string
}

func (c *GameClient) login() error {
	b, err := json.Marshal(types.LoginData{
		ClientID: c.clientID,
		Username: c.Username,
	})
	if err != nil {
		return err
	}
	msg := types.Message{
		Type: "login",
		Data: b,
	}
	return c.conn.WriteJSON(msg)
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

	for {
		x := rand.Intn(1000)
		y := rand.Intn(1000)
		state := types.PlayerState{
			Health:   100,
			Position: types.Position{X: x, Y: y},
		}
		bState, err := json.Marshal(state)
		if err != nil {
			log.Fatal(err)
		}
		msg := types.Message{
			Type: "playerState",
			Data: bState,
		}
		conn.WriteJSON(msg)
		time.Sleep(time.Millisecond * 100)
	}
}
