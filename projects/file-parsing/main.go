package main

import (
	"fmt"
	"os"

	jsonparser "github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/fileparser/json"
)

func getHighestScorePlayer(players jsonparser.Players) string {
	highestScore := players[0].HighScore
	highestScorePlayer := ""
	for _, player := range players {
		if player.HighScore >= highestScore {
			highestScore = player.HighScore
			highestScorePlayer = player.Name
		}
	}
	return highestScorePlayer
}

func getLowestScorePlayer(players jsonparser.Players) string {
	lowestScore := players[0].HighScore
	lowestScorePlayer := ""
	for _, player := range players {
		if player.HighScore <= lowestScore {
			lowestScore = player.HighScore
			lowestScorePlayer = player.Name
		}
	}
	return lowestScorePlayer
}

func main() {
	data, err := jsonparser.ParseJSON("./examples/json.txt")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	highestScorePlayer := getHighestScorePlayer(data)
	lowestScorePlayer := getLowestScorePlayer(data)
	fmt.Printf("The player with the highest score is '%s'\n", highestScorePlayer)
	fmt.Printf("The player with the lowest score is '%s'\n", lowestScorePlayer)
}
