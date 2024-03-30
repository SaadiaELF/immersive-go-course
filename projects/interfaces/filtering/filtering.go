package filtering

import (
	"io"
)

type filteringPipe struct {
	fw io.Writer
}

func NewFilteringPipe(w io.Writer) io.Writer {
	return &filteringPipe{fw: w}
}

func (f filteringPipe) Write(p []byte) (n int, err error) {
	filtered := []byte{}
	for _, value := range p {
		if value < '0' || value > '9' {
			filtered = append(filtered, value)
		}
	}
	n, err = f.fw.Write(filtered)
	if err != nil {
		return n, err
	}

	return n, nil
}
