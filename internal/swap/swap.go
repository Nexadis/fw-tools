package swap

import (
	"github.com/Nexadis/fw-tools/internal/config"
)

type Swapper struct {
	Input  string
	Output string
	Config config.Swap
}

func (s Swapper) Swap() {

}

func swapBits(b uint8) uint8 {
	var o uint8
	for i := 0; i < 8; i++ {
		o += ((b & (1 << i)) >> i) << (7 - i)
	}
	return o

}

func swapHalf(h uint8) uint8 {
	top := (h & 0xF0) >> 4
	bot := (h & 0x0F) << 4
	return top + bot
}

func swapWord(w uint16) uint16 {
	top := (w & 0xFF00) >> 8
	bot := (w & 0x00FF) << 8
	return top + bot
}

func swapDWord(d uint32) uint32 {
	top := (d & 0xFFFF0000) >> 16
	bot := (d & 0x0000FFFF) << 16
	return top + bot
}
