package xlsx

import (
	"bytes"
)

type Buffer struct {
	bytes.Buffer
}

func (b *Buffer) Get() []byte {
	return b.Buffer.Bytes()
}
