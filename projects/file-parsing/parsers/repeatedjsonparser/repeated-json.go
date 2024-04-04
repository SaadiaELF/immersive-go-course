package repeatedjsonparser

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/types"
)

type Parser struct{}

func (p Parser) Parse(filename string) (players types.Players, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		player := types.Player{}
		err = json.Unmarshal([]byte(line), &player)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}
