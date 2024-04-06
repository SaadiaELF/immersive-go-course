package jsonparser

import (
	"testing"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/utils"
	"github.com/stretchr/testify/require"
)

func TestJSONParser(t *testing.T) {
	myParser := Parser{}

	actual, err := myParser.Parse("../../examples/json.txt")
	require.NoError(t, err)
	expected := utils.Players
	require.Equal(t, expected, actual)

}
