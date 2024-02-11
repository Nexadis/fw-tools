package swap

import (
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"runtime"
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
			err := s.swap(context.TODO(), tn.prepare(), buf)
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

var Err error

const dataSize = 1 << 20

var filesCount int = runtime.GOMAXPROCS(0)

type bench struct {
	name string
	conf config.Swap
	size int
}

func BenchmarkBits(b *testing.B) {
	benches := []bench{
		{
			"bits",
			config.Swap{Bits: true},
			dataSize,
		},
	}
	for _, tb := range benches {
		b.Run(tb.name, func(b *testing.B) {
			s := New(tb.conf)
			ctx := context.TODO()
			inbuf := make([]byte, tb.size)
			rand.Read(inbuf)
			outbuf := make([]byte, tb.size)
			b.ResetTimer()
			var err error
			for i := 0; i < b.N; i++ {
				r := bytes.NewBuffer(inbuf)
				w := bytes.NewBuffer(outbuf)
				err = s.swap(ctx, r, w)
			}
			Err = err
		})
	}

}
func BenchmarkHalfs(b *testing.B) {
	benches := []bench{
		{
			"halfs",
			config.Swap{Halfs: true},
			dataSize,
		},
	}
	for _, tb := range benches {
		b.Run(tb.name, func(b *testing.B) {
			s := New(tb.conf)
			ctx := context.TODO()
			inbuf := make([]byte, tb.size)
			rand.Read(inbuf)
			outbuf := make([]byte, tb.size)
			b.ResetTimer()
			var err error
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r := bytes.NewBuffer(inbuf)
				w := bytes.NewBuffer(outbuf)
				b.StartTimer()
				err = s.swap(ctx, r, w)
			}
			Err = err
		})
	}

}
func BenchmarkBytes(b *testing.B) {
	benches := []bench{
		{
			"bytes",
			config.Swap{Bytes: true},
			dataSize,
		},
	}
	for _, tb := range benches {
		b.Run(tb.name, func(b *testing.B) {
			s := New(tb.conf)
			ctx := context.TODO()
			inbuf := make([]byte, tb.size)
			rand.Read(inbuf)
			outbuf := make([]byte, tb.size)
			b.ResetTimer()
			var err error
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r := bytes.NewBuffer(inbuf)
				w := bytes.NewBuffer(outbuf)
				b.StartTimer()
				err = s.swap(ctx, r, w)
			}
			Err = err
		})
	}

}
func BenchmarkWords(b *testing.B) {
	benches := []bench{
		{
			"words",
			config.Swap{Words: true},
			dataSize,
		},
	}
	for _, tb := range benches {
		b.Run(tb.name, func(b *testing.B) {
			s := New(tb.conf)
			ctx := context.TODO()
			inbuf := make([]byte, tb.size)
			rand.Read(inbuf)
			outbuf := make([]byte, tb.size)
			b.ResetTimer()
			var err error
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r := bytes.NewBuffer(inbuf)
				w := bytes.NewBuffer(outbuf)
				b.StartTimer()
				err = s.swap(ctx, r, w)
			}
			Err = err
		})
	}

}
func BenchmarkDwords(b *testing.B) {
	benches := []bench{
		{
			"dwords",
			config.Swap{Dwords: true},
			dataSize,
		},
	}
	for _, tb := range benches {
		b.Run(tb.name, func(b *testing.B) {
			s := New(tb.conf)
			ctx := context.TODO()
			inbuf := make([]byte, tb.size)
			rand.Read(inbuf)
			outbuf := make([]byte, 0, tb.size)
			b.ResetTimer()
			var err error
			for i := 0; i < b.N; i++ {
				r := bytes.NewBuffer(inbuf)
				w := bytes.NewBuffer(outbuf)
				err = s.swap(ctx, r, w)
			}
			Err = err
		})
	}

}

func BenchmarkRun(b *testing.B) {
	dir := os.TempDir()
	benches := []bench{
		{
			"bits",
			config.Swap{Bits: true},
			dataSize,
		},
		{
			"bytes",
			config.Swap{Bytes: true},
			dataSize,
		},
	}
	names := make([]string, 0, filesCount)
	for i := 0; i < filesCount; i++ {
		names = append(names, fmt.Sprintf("%s/test-fw-tools_%d.bin", dir, i))
	}
	for _, tb := range benches {
		b.Run(tb.name, func(b *testing.B) {
			buf := make([]byte, tb.size)
			rand.Read(buf)
			for _, name := range names {
				err := os.WriteFile(name, buf, 0777)
				if err != nil {
					b.Fatal(err)
				}
			}
			s := New(tb.conf)
			ctx := context.TODO()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s.Open(names)
				s.Run(ctx)
				s.Close()
			}

		})

	}
}
