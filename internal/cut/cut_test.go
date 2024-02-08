package cut

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nexadis/fw-tools/internal/config"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg config.Cut
	}
	tests := []struct {
		name string
		args args
		want *Cutter
	}{
		{
			"Simple create new Cutter",
			args{},
			&Cutter{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCutter_Open(t *testing.T) {
	type args struct {
		inputs []string
	}
	tests := []struct {
		name    string
		c       *Cutter
		args    args
		wantErr bool
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Open(tt.args.inputs); (err != nil) != tt.wantErr {
				t.Errorf("Cutter.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCutter_Close(t *testing.T) {
	tests := []struct {
		name    string
		c       *Cutter
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Close(); (err != nil) != tt.wantErr {
				t.Errorf("Cutter.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCutter_Run(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		prepare func() *Cutter
		want    []byte
		wantErr bool
	}{
		{
			"Cut pages",
			func() *Cutter {
				c := &Cutter{
					Config: config.Cut{
						PageSize: 8,
						SkipSize: 4,
					},

					inputs: []io.ReadCloser{io.NopCloser(bytes.NewBufferString("aaaabbbbccccddddeeeeffff"))},
				}
				return c
			},
			[]byte("aaaabbbbddddeeee"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.NotNil(t, tt.prepare)
			c := tt.prepare()
			buf := &bytes.Buffer{}
			c.outputs = []io.WriteCloser{NopWCloser(buf)}
			if err := c.Run(context.TODO()); (err != nil) != tt.wantErr {
				t.Errorf("Cutter.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.want, buf.Bytes())

		})
	}
}

type nopCloser struct {
	io.Writer
}

func (nc *nopCloser) Close() error {
	return nil
}

func NopWCloser(w io.Writer) io.WriteCloser {
	return &nopCloser{w}

}
