package main

import (
	"fmt"
	"os"

	jsonparser "github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/fileparser"
)

func getHighScorePlayer(players jsonparser.Players) string {
	highScore := 0
	highScorePlayer := ""
	for _, player := range players {
		if player.HighScore > highScore {
			highScore = player.HighScore
			highScorePlayer = player.Name
		}
	}
	return highScorePlayer
}

func main() {
	data, err := jsonparser.ParseJSON("./examples/json.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	highScorePlayer := getHighScorePlayer(data)
	fmt.Printf("The player with the highest score is '%s'\n", highScorePlayer)
}
