package OurBuffer

// import "bytes"

type OurByteBuffer struct {
	bytes []byte
}

func NewBuffer(b []byte) *OurByteBuffer {
	return &OurByteBuffer{bytes: b}
}

func (b *OurByteBuffer) Bytes() []byte {
	return b.bytes
}

func (b *OurByteBuffer) Write(p []byte) (n int) {
	b.bytes = append(b.bytes, p...)
	return len(p)
}

func (b *OurByteBuffer) Read(p []byte) (n int) {
	if len(p) > len(b.bytes) {
		b.bytes = b.bytes[len(b.bytes):]
		return len(b.bytes)
	} else {
		b.bytes = b.bytes[len(p):]
		return len(p)
	}
}

func (b *OurByteBuffer) String() string {
	return string(b.bytes)
}
