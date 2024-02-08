package swap

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nexadis/fw-tools/internal/config"
)

const (
	_ = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
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
		arg  []byte
		want []byte
	}{
		{
			name: "Simple swap dword",
			arg:  []byte{0xAB, 0xCD, 0x12, 0x34, 0x56, 0x78, 0x90, 0xEF},
			want: []byte{0x56, 0x78, 0x90, 0xEF, 0xAB, 0xCD, 0x12, 0x34},
		},
	}
	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			SwapDwords(tn.arg)
			assert.Equal(t, tn.want, tn.arg)
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
		{
			"Inverse words only",
			config.Swap{
				Words: true,
			},
			func() io.Reader {
				return bytes.NewReader(bytes.Repeat([]byte{0x01, 0x02, 0x03, 0x04}, 10))
			},
			bytes.Repeat([]byte{0x03, 0x04, 0x01, 0x02}, 10),
		},
		{
			"Inverse dwords only",
			config.Swap{
				Dwords: true,
			},
			func() io.Reader {
				return bytes.NewReader(bytes.Repeat([]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}, 10))
			},
			bytes.Repeat([]byte{0x05, 0x06, 0x07, 0x08, 0x01, 0x02, 0x03, 0x04}, 10),
		},
	}
	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			s := Swapper{
				Config: tn.conf,
			}
			buf := bytes.NewBuffer(make([]byte, 0, len(tn.want)))
			err := s.Swap(context.TODO(), tn.prepare(), buf)
			require.NoError(t, err)
			require.Equal(t, tn.want, buf.Bytes())

		})
	}
}

func TestCheckLen(t *testing.T) {
	tests := []struct {
		name    string
		conf    config.Swap
		len     int64
		wantErr bool
	}{
		{
			"No align",
			config.Swap{
				Bits:  true,
				Halfs: true,
			},
			123,
			false,
		},
		{
			"Invalid Align word",
			config.Swap{
				Bytes: true,
			},
			123,
			true,
		},
		{
			"Valid Align word",
			config.Swap{
				Bytes: true,
			},
			122,
			false,
		},
		{
			"Invalid Align dword",
			config.Swap{
				Words: true,
			},
			6,
			true,
		},
		{
			"Valid Align dword",
			config.Swap{
				Words: true,
			},
			8,
			false,
		},
		{
			"Invalid Align qword",
			config.Swap{
				Dwords: true,
			},
			12,
			true,
		},
		{
			"Valid Align qword",
			config.Swap{
				Dwords: true,
			},
			16,
			false,
		},
	}
	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			s := Swapper{
				Config: tn.conf,
			}
			err := s.checkLen(tn.len)
			if tn.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

		})
	}

}

func TestRun(t *testing.T) {
	tests := []struct {
		name    string
		conf    config.Swap
		in      []byte
		out     []byte
		wantErr bool
	}{
		{
			"Test byte file",
			config.Swap{
				Bits:  true,
				Halfs: true,
			},
			bytes.Repeat([]byte{0b1010_1011}, 0x40000),
			bytes.Repeat([]byte{0b0101_1101}, 0x40000),
			false,
		},
		{
			"Test word and dword",
			config.Swap{
				Bytes: true,
				Words: true,
			},
			bytes.Repeat([]byte{0xAD, 0xEF, 0x01, 0x23}, 0x40000),
			bytes.Repeat([]byte{0x23, 0x01, 0xEF, 0xAD}, 0x40000),
			false,
		},
	}
	inname := os.TempDir() + "/test_in.bin"
	os.Remove(inname)

	for _, tn := range tests {
		t.Run(tn.name, func(t *testing.T) {
			defer os.Remove(inname)

			s := Swapper{
				Config: tn.conf,
			}
			outname := s.outName(inname)
			defer os.Remove(outname)
			err := os.WriteFile(inname, tn.in, 0777)
			require.NoError(t, err)
			require.NoError(t, s.Open([]string{inname}))

			err = s.Run(context.Background())
			if tn.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.NoError(t, s.Close())
			data, err := os.ReadFile(outname)
			require.NoError(t, err)
			assert.Equal(t, tn.out, data)

		})
	}

}

func BenchmarkRun(b *testing.B) {
	tests := []struct {
		name    string
		conf    config.Swap
		in      []byte
		wantErr bool
	}{
		{
			"Bench byte file",
			config.Swap{
				Bits:  true,
				Halfs: true,
			},
			[]byte{0b1010_1011},
			false,
		},
		{
			"Bench word and dword file",
			config.Swap{
				Bytes: true,
				Words: true,
			},
			[]byte{0xAD},
			false,
		},
	}
	inname := os.TempDir() + "/test_in.bin"
	os.Remove(inname)

	for _, tt := range tests {
		sizes := []int{KiB, MiB, 10 * MiB, 100 * MiB, 1 * GiB}
		for _, size := range sizes {
			b.Run(fmt.Sprintf("%s_%d", tt.name, size), func(b *testing.B) {
				defer os.Remove(inname)

				f, _ := os.OpenFile(inname, os.O_CREATE|os.O_WRONLY, 0755)
				f.Write(bytes.Repeat(tt.in, size))
				f.Close()

				s := Swapper{
					Config: tt.conf,
				}
				outname := s.outName(inname)
				defer os.Remove(outname)

				s.Open([]string{inname})
				b.Log("end prepare")
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					s.Run(context.TODO())
				}
				s.Close()

			})
		}
	}

}

func BenchmarkSwapDwords(b *testing.B) {
	tests := []struct {
		name string
		size int
	}{
		{
			"Byte slice Swap 1MiB",
			1 * MiB,
		},
		{
			"Byte slice Swap 100MiB",
			100 * MiB,
		},
		{
			"Byte slice Swap 1GiB",
			1 * GiB,
		},
	}
	for _, tn := range tests {
		b.Run(tn.name, func(b *testing.B) {
			data := bytes.Repeat([]byte{0xFA}, tn.size)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				SwapDwords(data)
			}

		})
	}
	utests := []struct {
		name     string
		sizeFile int
		bufSize  int
	}{
		{
			"Uint Swap 1MiB",
			1 * MiB,
			8,
		},
		{
			"Uint Swap 100MiB",
			100 * MiB,
			8,
		},
		{
			"Uint Swap 1GiB buf=8",
			1 * GiB,
			0x8,
		},
		{
			"Uint Swap 1GiB buf=0x100",
			1 * GiB,
			0x100,
		},
		{
			"Uint Swap 1GiB buf=0x400",
			1 * GiB,
			0x400,
		},
		{
			"Uint Swap 1GiB buf=0x1000",
			1 * GiB,
			0x1000,
		},
	}
	for _, tn := range utests {
		b.Run(tn.name, func(b *testing.B) {
			data := bytes.Repeat([]byte{0xFA}, tn.sizeFile)
			r := bytes.NewReader(data)
			buf := make([]byte, tn.bufSize)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for n, err := r.Read(buf); n != 0; n, err = r.Read(buf) {
					if err != nil && err != io.EOF {
						return
					}
					var w uint64
					for i := 0; i < len(buf); i += 8 {
						w = binary.BigEndian.Uint64(buf[i : i+8])
						w = SwapUInt(w)
						binary.BigEndian.PutUint64(buf[i:i+8], w)
					}
				}
			}
		})
	}
}
