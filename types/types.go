package types

type LoginData struct {
	ClientID int    `json:"client_id"`
	Username string `json:"username"`
}

type Message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type PlayerState struct {
	Health   int      `json:"health"`
	Position Position `json:"position"`
}
