package utils

import (
	"testing"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/types"
	"github.com/stretchr/testify/require"
)

var players = Players
func TestGetHighestScorePlayer(t *testing.T) {
	actual := GetHighestScorePlayer(players)
	expected := types.Player{Name: "Prisha", HighScore: 30}
	require.Equal(t, expected, actual)
}

func TestGetLowestScorePlayer(t *testing.T) {
	actual := GetLowestScorePlayer(players)
	expected := types.Player{Name: "Charlie", HighScore: -1}
	require.Equal(t, expected, actual)
}
