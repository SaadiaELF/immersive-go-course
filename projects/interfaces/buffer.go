package OurBuffer

import "io"

type OurByteBuffer struct {
	bytes []byte
}

func NewBuffer(b []byte) *OurByteBuffer {
	return &OurByteBuffer{bytes: b}
}

func (b *OurByteBuffer) Bytes() []byte {
	return b.bytes
}

func (b *OurByteBuffer) Write(p []byte) (n int, err error) {
	b.bytes = append(b.bytes, p...)
	return len(p), nil
}

func (b *OurByteBuffer) Read(p []byte) (n int, err error) {
	n = len(p)
	if n > len(b.bytes) {
		n = len(b.bytes)
	}

	copy(p, b.bytes[:n])
	
	b.bytes = b.bytes[n:]

	if len(b.bytes) == 0 {
		err = io.EOF
	}

	return n, err
}

func (b *OurByteBuffer) String() string {
	return string(b.bytes)
}
