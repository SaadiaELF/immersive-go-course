package main

import (
	"testing"

	jsonparser "github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/fileparser"
	"github.com/stretchr/testify/require"
)

func TestGetHighScorePlayer(t *testing.T) {
	players := jsonparser.Players{
		{Name: "Aya", HighScore: 10},
		{Name: "Prisha", HighScore: 30},
		{Name: "Charlie", HighScore: -1},
		{Name: "Margot", HighScore: 25},
	}
	actual := getHighScorePlayer(players)
	expected := "Prisha"
	require.Equal(t, expected, actual)
}
