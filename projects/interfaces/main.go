package main

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

func main() {

}
