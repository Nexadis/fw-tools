package cut

import (
	"context"
	"io"
	"os"
	"strings"

	"github.com/Nexadis/fw-tools/internal/config"
)

type Cutter struct {
	input  io.ReadCloser
	output io.WriteCloser
	Config config.Cut
}

func New(cfg config.Cut) *Cutter {
	return &Cutter{
		Config: cfg,
	}
}

func (c *Cutter) Open(input string, output string) error {
	fi, err := os.Open(input)
	if err != nil {
		return err
	}
	c.input = fi
	if output == "" {
		name, _ := strings.CutSuffix(input, ".bin")
		output = name + "-cutted.bin"
	}
	fo, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0766)
	if err != nil {
		return err
	}
	c.output = fo
	return nil
}

func (c *Cutter) Close() error {
	err := c.input.Close()
	if err != nil {
		return err
	}
	return c.output.Close()

}

func (c *Cutter) Run(ctx context.Context) error {
	skip := io.Discard
	for {
		n, err := io.CopyN(c.output, c.input, int64(c.Config.Page))
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		io.CopyN(skip, c.input, int64(c.Config.Skip))
	}

}