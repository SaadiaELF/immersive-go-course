package OurBuffer

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
	if len(p) >= len(b.bytes) {
		b.bytes = b.bytes[len(b.bytes):]
		copy(p, b.bytes)
	}
	if len(p) < len(b.bytes) {
		b.bytes = b.bytes[len(p):]
	}
	return len(p), nil
}

func (b *OurByteBuffer) String() string {
	return string(b.bytes)
}
