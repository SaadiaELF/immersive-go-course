package repeatedjsonparser

import (
	"testing"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/utils"
	"github.com/stretchr/testify/require"
)

func TestRepeatedJSONParser(t *testing.T) {
	myParser := Parser{}
	actual, err := myParser.Parse("../../examples/repeated-json.txt")
	require.NoError(t, err)
	expected := utils.Players
	require.Equal(t, expected, actual)
}
