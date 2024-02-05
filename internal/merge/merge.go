package merge

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Nexadis/fw-tools/internal/config"
)

type Merger struct {
	inputs []io.ReadCloser
	output io.WriteCloser
	Config config.Merge
}

func New(cfg config.Merge) *Merger {
	return &Merger{
		Config: cfg,
	}
}

func (m *Merger) Open(inputs []string, output string) error {
	// size - size of file, must be the same for all files
	var size int64 = -1
	m.inputs = make([]io.ReadCloser, 0, len(inputs))
	for _, i := range inputs {
		in, err := os.Open(i)
		if err != nil {
			return fmt.Errorf("can't open file for merging: %w", err)
		}
		s, err := in.Stat()
		if err != nil {
			return fmt.Errorf("can't get file stat for merging: %w", err)
		}

		// save size of first file
		if size == -1 {
			size = s.Size()
		}

		// file with alternative size
		if size != s.Size() {
			return fmt.Errorf("size of file is not same, problem with:%s", s.Name())
		}
		m.inputs = append(m.inputs, in)
	}
	if output == "" {
		output = "merged.bin"
	}
	o, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0766)
	m.output = o
	return err
}

func (m *Merger) Close() error {
	for _, i := range m.inputs {
		err := i.Close()
		if err != nil {
			return err
		}
	}
	return m.output.Close()
}

func (m *Merger) Run(ctx context.Context) error {
	switch {
	case m.Config.ByBit:
		return m.bits(ctx)
	case m.Config.ByByte:
		return m.bytes(ctx, 1)
	case m.Config.ByWord:
		return m.bytes(ctx, 2)
	case m.Config.ByDword:
		return m.bytes(ctx, 4)
	default:
		return errors.New("unexpected mode, choose one")
	}
}

func (m *Merger) bytes(ctx context.Context, size int) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		for _, i := range m.inputs {
			n, err := io.CopyN(m.output, i, int64(size))
			if err != nil && err != io.EOF {
				return err
			}
			// first empty reader, all the same size
			if n == 0 {
				return nil
			}
		}
	}

}

func (m *Merger) bits(ctx context.Context) error {
	bPerInput := make([]byte, len(m.inputs))
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		for i, r := range m.inputs {
			n, err := r.Read(bPerInput[i : i+1])
			if err != nil && err != io.EOF {
				return err
			}
			if n == 0 {
				return nil
			}
		}
		out := make([]byte, len(m.inputs))
		// bitOffOut - bit offset in output sequence of bytes
		bitOffOut := 0
		for bitOff := 0; bitOff < 8; bitOff++ {
			mask := byte(1 << bitOff)
			for _, b := range bPerInput {
				bv := ((b & mask) >> bitOff) << byte(bitOffOut%8)
				out[bitOffOut/8] |= bv
				bitOffOut++
			}
		}
		m.output.Write(out)
	}

}
