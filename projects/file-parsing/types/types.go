package types

type Player struct {
	Name      string `json:"name"`
	HighScore int32  `json:"high_score"`
}

type Players []Player

type Parser interface {
	Parse(string) (Players, error)
}
