package utils

import (
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/types"
)

var Players = types.Players{
	{Name: "Aya", HighScore: 10},
	{Name: "Prisha", HighScore: 30},
	{Name: "Charlie", HighScore: -1},
	{Name: "Margot", HighScore: 25},
}

func GetHighestScorePlayer(players types.Players) types.Player {
	highestScorePlayer := players[0]
	for _, player := range players {
		if player.HighScore >= highestScorePlayer.HighScore {
			highestScorePlayer.HighScore = player.HighScore
			highestScorePlayer.Name = player.Name
		}
	}
	return highestScorePlayer
}

func GetLowestScorePlayer(players types.Players) types.Player {
	lowestScorePlayer := players[0]
	for _, player := range players {
		if player.HighScore <= lowestScorePlayer.HighScore {
			lowestScorePlayer.HighScore = player.HighScore
			lowestScorePlayer.Name = player.Name
		}
	}
	return lowestScorePlayer
}
