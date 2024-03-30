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
	n, err = f.fw.Write(p)
	if err != nil {
		return n, err
	}

	return n, nil
}
