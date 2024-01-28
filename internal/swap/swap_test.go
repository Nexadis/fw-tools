package swap

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nexadis/fw-tools/internal/config"
)

func TestInverseBits(t *testing.T) {
	tests := []struct {
		name string
		arg  uint8
		want uint8
	}{
		{
			name: "Simple inverse bits",
			arg:  0b10111001,
			want: 0b10011101,
		},
	}
	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			assert.Equal(t, tn.want, InverseBits(tn.arg))
		})
	}
}

func TestSwapHalf(t *testing.T) {
	tests := []struct {
		name string
		arg  uint8
		want uint8
	}{
		{
			name: "Simple swap half",
			arg:  0b10111001,
			want: 0b10011011,
		},
	}
	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			assert.Equal(t, tn.want, SwapHalf(tn.arg))
		})
	}
}

func TestSwapBytes(t *testing.T) {
	tests := []struct {
		name string
		arg  uint16
		want uint16
	}{
		{
			name: "Simple swap bytes",
			arg:  0xABCD,
			want: 0xCDAB,
		},
	}
	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			assert.Equal(t, tn.want, SwapBytes(tn.arg))
		})
	}
}

func TestSwapWord(t *testing.T) {
	tests := []struct {
		name string
		arg  uint32
		want uint32
	}{
		{
			name: "Simple swap word",
			arg:  0xABCD1234,
			want: 0x1234ABCD,
		},
	}
	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			assert.Equal(t, tn.want, SwapWords(tn.arg))
		})
	}
}

func TestSwapDWord(t *testing.T) {
	tests := []struct {
		name string
		arg  uint64
		want uint64
	}{
		{
			name: "Simple swap dword",
			arg:  0xABCD1234567890EF,
			want: 0x567890EFABCD1234,
		},
	}
	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			assert.Equal(t, tn.want, SwapDwords(tn.arg))
		})
	}
}

func TestSwap(t *testing.T) {
	tests := []struct {
		name    string
		conf    config.Swap
		prepare func() io.Reader
		want    []byte
	}{
		{
			"Inverse bits only",
			config.Swap{
				Bits: true,
			},
			func() io.Reader {
				return bytes.NewReader(bytes.Repeat([]byte{0b1100_1101, 0b1010_1011, 0b0101_0001}, 10))
			},
			bytes.Repeat([]byte{0b1011_0011, 0b1101_0101, 0b1000_1010}, 10),
		},
		{
			"Inverse half only",
			config.Swap{
				Halfs: true,
			},
			func() io.Reader {
				return bytes.NewReader(bytes.Repeat([]byte{0b1100_1101, 0b1010_1011, 0b0101_0001}, 10))
			},
			bytes.Repeat([]byte{0b1101_1100, 0b1011_1010, 0b0001_0101}, 10),
		},
		{
			"Inverse bytes only",
			config.Swap{
				Bytes: true,
			},
			func() io.Reader {
				return bytes.NewReader(bytes.Repeat([]byte{0x01, 0x02, 0x03, 0x04}, 10))
			},
			bytes.Repeat([]byte{0x02, 0x01, 0x04, 0x03}, 10),
		},
	}
	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			s := Swapper{
				Config: tn.conf,
			}
			buf := bytes.NewBuffer(make([]byte, 0, len(tn.want)))
			err := s.Swap(tn.prepare(), buf)
			require.NoError(t, err)
			require.Equal(t, tn.want, buf.Bytes())

		})
	}
}
