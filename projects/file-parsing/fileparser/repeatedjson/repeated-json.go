package repeatedjson

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"
)

type Player struct {
	Name      string `json:"name"`
	HighScore int    `json:"high_score"`
}

type Players []Player

func RepeatedJSONParser(filename string) (players Players, err error) {
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
		player := Player{}
		err = json.Unmarshal([]byte(line), &player)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}
	return players, nil
}
