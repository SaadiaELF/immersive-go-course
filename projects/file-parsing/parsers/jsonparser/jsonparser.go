package jsonparser

import (
	"encoding/json"
	"io"
	"os"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/types"
)

type Parser struct{}

func (p Parser) Parse(filename string) (players types.Players, err error) {
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
