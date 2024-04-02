package main

import (
	"testing"

	jsonparser "github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/fileparser"
	"github.com/stretchr/testify/require"
)

var players = jsonparser.Players{
	{Name: "Aya", HighScore: 10},
	{Name: "Prisha", HighScore: 30},
	{Name: "Charlie", HighScore: -1},
	{Name: "Margot", HighScore: 25},
}

func TestGetHighestScorePlayer(t *testing.T) {
	actual := getHighestScorePlayer(players)
	expected := "Prisha"
	require.Equal(t, expected, actual)
}

func TestGetLowestScorePlayer(t *testing.T) {
	actual := getLowestScorePlayer(players)
	expected := "Charlie"
	require.Equal(t, expected, actual)
}
