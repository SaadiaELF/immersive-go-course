package main

import (
	"fmt"
	"os"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/binary"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/jsonparser"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/parsers/repeatedjsonparser"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/types"
	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/utils"
)

func main() {

	// open a directory and read all files
	// for each file, parse it and get the highest and lowest score player
	// print the highest and lowest score player
	dataMap := make(map[string]types.Players)
	dir := "./examples"
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	var myParser types.Parser
	var data types.Players
	for _, file := range files {
		filepath := dir + "/" + file.Name()

		switch filepath {
		case "./examples/json.txt":
			myParser = jsonparser.Parser{}
		case "./examples/repeated-json.txt":
			myParser = repeatedjsonparser.Parser{}
		case "./examples/custom-binary-be.bin":
			myParser = binary.Parser{}
		case "./examples/custom-binary-le.bin":
			myParser = binary.Parser{}
		default:
			fmt.Fprintf(os.Stderr, "Unknown file type: %s\n", filepath)
		}
		data, err = myParser.Parse(filepath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing file %s: %v\n", filepath, err)
			continue // Skip to the next file if parsing fails
		}

		dataMap[filepath] = data
	}

	for filepath, data := range dataMap {
		highestScorePlayer := utils.GetHighestScorePlayer(data)
		lowestScorePlayer := utils.GetLowestScorePlayer(data)

		fmt.Printf("The player with the highest score in %s is %s with %d\n", filepath, highestScorePlayer.Name, highestScorePlayer.HighScore)
		fmt.Printf("The player with the lowest score in %s is %s with %d\n", filepath, lowestScorePlayer.Name, lowestScorePlayer.HighScore)
	}
}
