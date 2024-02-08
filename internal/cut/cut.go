package cut

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"golang.org/x/sync/errgroup"

	"github.com/Nexadis/fw-tools/internal/config"
)

type Cutter struct {
	inputs  []io.ReadCloser
	outputs []io.WriteCloser
	Config  config.Cut
}

func New(cfg config.Cut) *Cutter {
	return &Cutter{
		Config: cfg,
	}
}

func (c *Cutter) Open(inputs []string) error {
	for _, in := range inputs {
		fi, err := os.Open(in)
		if err != nil {
			return err
		}
		c.inputs = append(c.inputs, fi)
		name, _ := strings.CutSuffix(in, ".bin")
		output := name + "-cutted.bin"
		fo, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			return err
		}
		c.outputs = append(c.outputs, fo)
	}
	return nil
}

func (c *Cutter) Close() error {
	var err error
	for _, in := range c.inputs {
		err = errors.Join(in.Close(), err)
	}
	for _, out := range c.outputs {
		err = errors.Join(out.Close(), err)
	}
	return err
}

func (c *Cutter) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(16)
	for i := 0; i < len(c.inputs); i++ {
		i := i
		grp.Go(func() error {
			return c.Cut(ctx, c.inputs[i], c.outputs[i])
		})

	}
	return grp.Wait()

}

func (c *Cutter) Cut(ctx context.Context, i io.Reader, o io.Writer) error {
	skip := io.Discard
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		n, err := io.CopyN(o, i, int64(c.Config.PageSize))
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			return nil
		}
		io.CopyN(skip, i, int64(c.Config.SkipSize))
	}

}
