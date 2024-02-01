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
	var size int64
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
		if size == 0 {
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
		return m.bytes(ctx)
	case m.Config.ByWord:
		return m.words(ctx)
	case m.Config.ByDword:
		return m.dwords(ctx)
	default:
		panic("unexpected mode")
	}
}

func (m *Merger) bits(ctx context.Context) error {
	return nil
}

func (m *Merger) bytes(ctx context.Context) error {
	b := []byte{0}
	for _, i := range m.inputs {
		n, err := i.Read(b)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			return nil
		}
		m.output.Write(b)
	}
	return nil
}

func (m *Merger) words(ctx context.Context) error {
	b := []byte{0, 0}
	for _, i := range m.inputs {
		n, err := i.Read(b)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			return nil
		}
		m.output.Write(b)
	}
	return nil
}

func (m *Merger) dwords(ctx context.Context) error {
	b := []byte{0, 0, 0, 0}
	for _, i := range m.inputs {
		n, err := i.Read(b)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			return nil
		}
		m.output.Write(b)
	}
	return nil
}
