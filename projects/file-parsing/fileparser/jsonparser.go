package jsonparser

import (
	"encoding/json"
	"io"
	"os"
)

type Player struct {
	Name      string `json:"name"`
	HighScore int    `json:"high_score"`
}

type Players []Player

func ParseJSON(filename string) (players Players, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, &players)
	return players, err
}
