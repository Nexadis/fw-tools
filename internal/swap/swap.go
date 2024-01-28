package swap

import (
	"bufio"
	"io"
	"os"

	"github.com/Nexadis/fw-tools/internal/config"
)

type Swapper struct {
	Input  string
	Output string
	Config config.Swap
}

func (s Swapper) Swap(i io.Reader, o io.Writer) error {
	buf := make([]byte, 8)
	for n, err := i.Read(buf); n != 0; n, err = i.Read(buf) {
		if err != nil && err != io.EOF {
			return err
		}
		if s.Config.Bits {
			for i, b := range buf {
				buf[i] = InverseBits(uint8(b))
			}
		}
		_, err := o.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Swapper) Run() error {
	in, err := os.OpenFile(s.Input, os.O_RDONLY, 0755)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(s.Output, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}
	defer out.Close()
	bufin := bufio.NewReader(in)
	bufout := bufio.NewWriter(out)
	return s.Swap(bufin, bufout)
}

func InverseBits(b uint8) uint8 {
	var o uint8
	for i := 0; i < 8; i++ {
		o += ((b & (1 << i)) >> i) << (7 - i)
	}
	return o

}

func SwapHalf(b uint8) uint8 {
	top := (b & 0xF0) >> 4
	bot := (b & 0x0F) << 4
	return top + bot
}

func SwapBytes(w uint16) uint16 {
	top := (w & 0xFF00) >> 8
	bot := (w & 0x00FF) << 8
	return top + bot
}

func SwapWords(d uint32) uint32 {
	top := (d & 0xFFFF0000) >> 16
	bot := (d & 0x0000FFFF) << 16
	return top + bot
}

func SwapDwords(q uint64) uint64 {
	top := (q & 0xFFFFFFFF00000000) >> 32
	bot := (q & 0x00000000FFFFFFFF) << 32
	return top + bot
}
