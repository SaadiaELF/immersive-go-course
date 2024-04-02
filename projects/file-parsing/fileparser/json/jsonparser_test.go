package jsonparser

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJSONParser(t *testing.T) {
	actual, err := JSONParser("../../examples/json.txt")
	require.NoError(t, err)
	expected := Players{{Name: "Aya", HighScore: 10}, {Name: "Prisha", HighScore: 30}, {Name: "Charlie", HighScore: -1}, {Name: "Margot", HighScore: 25}}
	require.Equal(t, expected, actual)

}
