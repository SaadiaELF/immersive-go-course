package binary

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/CodeYourFuture/immersive-go-course/projects/file-parsing/types"
)

type Parser struct{}

func byteOrder(file *os.File) (binary.ByteOrder, error) {
	var byteOrder binary.ByteOrder
	var endianBytes [2]byte
	_, err := file.Read(endianBytes[:])
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	if endianBytes == [2]byte{0xFE, 0xFF} {
		byteOrder = binary.BigEndian
	} else if endianBytes == [2]byte{0xFF, 0xFE} {
		byteOrder = binary.LittleEndian
	} else {
		return nil, fmt.Errorf("unknown byte order: %v", endianBytes)
	}

	return byteOrder, nil
}

func (p Parser) Parse(filename string) (players types.Players, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteOrder, err := byteOrder(file)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(file)

	for {
		// read the 4 bytes to get the score of the player
		var score int32

		err = binary.Read(r, byteOrder, &score)

		if err != nil {
			if err == io.EOF {
				break // exit the loop when reaching the end of the file
			}
			return nil, fmt.Errorf("error reading score: %v", err)
		}

		// read next bytes to get the name of the player till the null byte
		name, err := r.ReadString(0)
		name = name[:len(name)-1] // remove the null byte
		if err != nil {
			return nil, fmt.Errorf("error reading name: %v", err)
		}

		players = append(players, types.Player{Name: name, HighScore: score})
	}
	return players, nil
}
