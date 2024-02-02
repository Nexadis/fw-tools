package merge

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Nexadis/fw-tools/internal/config"
)

type Merger struct {
	inputs []io.Reader
	output io.Writer
	Config config.Merge
}

func New(cfg config.Merge) *Merger {
	return &Merger{
		Config: cfg,
	}
}

func (m *Merger) Open(inputs []string, output string) error {
	var size int64 = -1
	m.inputs = make([]io.Reader, 0, len(inputs))
	for _, i := range inputs {
		in, err := os.Open(i)
		if err != nil {
			return fmt.Errorf("can't open file for merging: %w", err)
		}
		s, err := in.Stat()
		if err != nil {
			return fmt.Errorf("can't open file for merging: %w", err)
		}
		if size == -1 {
			size = s.Size()
		}
		if size != s.Size() {
			return fmt.Errorf("size of file is not same, problem with:%s", s.Name())
		}
		m.inputs = append(m.inputs, bufio.NewReader(in))
	}
	o, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0766)
	m.output = bufio.NewWriter(o)
	return err
}

func (m *Merger) Run(ctx context.Context) error {
	switch {
	case m.Config.ByBit:
		return m.bits(ctx)
	case m.Config.ByByte:
		return m.mergeSize(ctx, 1)
	case m.Config.ByWord:
		return m.mergeSize(ctx, 2)
	case m.Config.ByDword:
		return m.mergeSize(ctx, 4)
	default:
		panic("unexpected mode")
	}
}

func (m *Merger) mergeSize(ctx context.Context, size int) error {
	b := make([]byte, size)
	var empties int
	for {
		if empties == len(m.inputs) {
			return nil
		}
		for _, i := range m.inputs {
			n, err := i.Read(b)
			if err != nil && err != io.EOF {
				return err
			}
			if n == 0 {
				empties += 1
				continue
			}
			m.output.Write(b)
		}
	}

}

func (m *Merger) bits(ctx context.Context) error {
	b := make([]byte, len(m.inputs))
	var empties int
	for {
		if empties == len(m.inputs) {
			return nil
		}
		for i, r := range m.inputs {
			n, err := r.Read(b[i : i+1])
			if err != nil && err != io.EOF {
				return err
			}
			if n == 0 {
				empties += 1
			}
		}
		out := make([]byte, len(m.inputs))
		num_bit := 0
		cur_byte := 0
		for bit := 0; bit < 8; bit++ {
			mask := byte(1 << bit)
			for _, file_byte := range b {
				bv := (file_byte & mask) << byte(num_bit)
				out[cur_byte] |= bv
				if num_bit == 7 {
					cur_byte += 1
					num_bit = 0
				}
			}
		}
		m.output.Write(out)
	}

}
