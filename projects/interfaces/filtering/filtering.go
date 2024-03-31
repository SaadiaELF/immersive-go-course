package filtering

import (
	"io"
)

type filteringPipe struct {
	fw io.Writer
}

func NewFilteringPipe(w io.Writer) io.Writer {
	return filteringPipe{fw: w}
}
func remove(slice []byte, s int) []byte {
	copy(slice[s:], slice[s+1:])
	return slice[:len(slice)-1]
}

func (f filteringPipe) Write(p []byte) (n int, err error) {
	for i, value := range p {
		if value < '0' || value > '9' {
			remove(p, i)
		}
	}
	n, err = f.fw.Write(p)
	if err != nil {
		return n, err
	}

	return n, nil
}
