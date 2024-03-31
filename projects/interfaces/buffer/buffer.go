package OurBuffer

import "io"

type OurByteBuffer struct {
	bytes []byte
}

// NewBuffer creates a new OurByteBuffer that takes ownership of the provided buffer.
// The caller should not modify or access the buffer after passing it to this function.
func NewBuffer(b []byte) *OurByteBuffer {
	return &OurByteBuffer{bytes: b}
}

// Bytes returns the actual underlying buffer of the OurByteBuffer.
// Modifying the returned buffer may lead to wrong read results.
func (b *OurByteBuffer) Bytes() []byte {
	return b.bytes
}

// Write appends the contents of p to the buffer, growing the buffer as needed.
func (b *OurByteBuffer) Write(p []byte) (n int, err error) {
	b.bytes = append(b.bytes, p...)
	return len(p), nil
}

// Read reads up to len(p) bytes into p from the buffer.
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

// String returns the contents of the buffer as a string.
func (b *OurByteBuffer) String() string {
	return string(b.bytes)
}
