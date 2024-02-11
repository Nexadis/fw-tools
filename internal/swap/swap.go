package swap

import (
	"bufio"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/Nexadis/fw-tools/internal/config"
)

var ErrAlign = errors.New("invalid align of input")

type Swapper struct {
	inputs  []io.ReadCloser
	outputs []io.WriteCloser
	Config  config.Swap
}

func New(cfg config.Swap) *Swapper {
	return &Swapper{
		Config: cfg,
	}
}

func (s *Swapper) Open(inputs []string) error {
	s.inputs = make([]io.ReadCloser, 0, len(inputs))
	s.outputs = make([]io.WriteCloser, 0, len(inputs))
	for _, i := range inputs {
		in, err := os.Open(i)
		if err != nil {
			return fmt.Errorf("can't open file '%s' for swapping: %w", i, err)
		}
		stat, err := in.Stat()
		if err != nil {
			return fmt.Errorf("can't get file stat '%s' for swapping: %w", i, err)
		}

		// file with alternative size
		err = s.checkLen(stat.Size())
		if err != nil {
			return fmt.Errorf("%w: %s", err, stat.Name())
		}
		s.inputs = append(s.inputs, in)
		o := s.outName(i)
		out, err := os.OpenFile(o, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return fmt.Errorf("can't create file '%s' for swapping: %w", o, err)
		}
		s.outputs = append(s.outputs, out)
	}
	return nil
}

func (s *Swapper) Close() error {
	var err error
	for _, in := range s.inputs {
		err = errors.Join(in.Close(), err)
	}
	for _, out := range s.outputs {
		err = errors.Join(out.Close(), err)
	}
	return err
}

func (s *Swapper) swap(ctx context.Context, i io.Reader, o io.Writer) error {
	buf := make([]byte, 0x400)
	for n, err := i.Read(buf); n != 0; n, err = i.Read(buf) {
		if err != nil && err != io.EOF {
			return err
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if s.Config.Bits {
			for i, b := range buf {
				buf[i] = InverseBits(b)
			}
		}

		if s.Config.Halfs {
			for i, b := range buf {
				buf[i] = SwapHalf(b)
			}
		}

		if s.Config.Bytes {
			var w uint16
			for i := 0; i < n; i += 2 {
				w = binary.BigEndian.Uint16(buf[i : i+2])
				w = SwapBytes(w)
				binary.BigEndian.PutUint16(buf[i:i+2], w)
			}
		}

		if s.Config.Words {
			var w uint32
			for i := 0; i < n; i += 4 {
				w = binary.BigEndian.Uint32(buf[i : i+4])
				w = SwapWords(w)
				binary.BigEndian.PutUint32(buf[i:i+4], w)
			}
		}

		if s.Config.Dwords {
			var w uint64
			for i := 0; i < n; i += 8 {
				w = binary.BigEndian.Uint64(buf[i : i+8])
				w = SwapUInt(w)
				binary.BigEndian.PutUint64(buf[i:i+8], w)
			}
		}

		_, err := o.Write(buf[:n])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s Swapper) checkLen(size int64) error {
	if s.Config.Bytes && size%2 != 0 {
		return fmt.Errorf("%w: should 2", ErrAlign)
	}
	if s.Config.Words && size%4 != 0 {
		return fmt.Errorf("%w: should 4", ErrAlign)
	}
	if s.Config.Dwords && size%8 != 0 {
		return fmt.Errorf("%w: should 8", ErrAlign)
	}
	return nil

}

func (s *Swapper) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	for i := 0; i < len(s.inputs); i++ {
		in := s.inputs[i]
		out := s.outputs[i]

		bufin := bufio.NewReader(in)
		bufout := bufio.NewWriter(out)
		grp.Go(func() error {
			defer bufout.Flush()
			return s.swap(ctx, bufin, bufout)
		})
	}
	return grp.Wait()
}

func InverseBits(b uint8) uint8 {
	var o uint8
	o |= (b & 0b1000_0000) >> 7
	o |= (b & 0b0100_0000) >> 5
	o |= (b & 0b0010_0000) >> 3
	o |= (b & 0b0001_0000) >> 1
	o |= (b & 0b0000_1000) << 1
	o |= (b & 0b0000_0100) << 3
	o |= (b & 0b0000_0010) << 5
	o |= (b & 0b0000_0001) << 7
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

func SwapDwords(b []byte) {
	beg := [4]byte{}
	copy(beg[:], b[:4])
	copy(b[:4], b[4:8])
	copy(b[4:8], beg[:])
}

func SwapUInt(d uint64) uint64 {
	top := (d & 0xFFFFFFFF00000000) >> 32
	bot := (d & 0x00000000FFFFFFFF) << 32
	return top + bot

}

func (s Swapper) outName(inName string) string {
	name, _ := strings.CutSuffix(inName, ".bin")
	if s.Config.Bits {
		name += "-bits"
	}
	if s.Config.Halfs {
		name += "-halfs"
	}
	if s.Config.Bytes {
		name += "-bytes"
	}
	if s.Config.Words {
		name += "-words"
	}
	if s.Config.Dwords {
		name += "-dwords"
	}
	name += ".bin"
	return name

}
