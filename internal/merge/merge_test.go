package merge

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Nexadis/fw-tools/internal/config"
)

func TestNew(t *testing.T) {
	type args struct {
		cfg config.Merge
	}
	tests := []struct {
		name string
		args args
		want *Merger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMerger_Open(t *testing.T) {
	type args struct {
		inputs []string
		output string
	}
	tests := []struct {
		name    string
		m       *Merger
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Open(tt.args.inputs, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("Merger.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMerger_Run(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		m       *Merger
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.m.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Merger.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMerger_mergeSize(t *testing.T) {
	type args struct {
		ctx  context.Context
		size int
	}
	tests := []struct {
		name    string
		prepare func() *Merger
		args    args
		out     *bytes.Buffer
		want    string
		wantErr bool
	}{
		{
			"Merge two files byte by byte",
			func() *Merger {
				m := &Merger{}
				m.inputs = []io.Reader{
					strings.NewReader("abcd"),
					strings.NewReader("1234"),
				}
				return m
			},
			args{
				nil,
				1,
			},
			bytes.NewBuffer(make([]byte, 0, 8)),
			"a1b2c3d4",
			false,
		},
		{
			"Merge three files byte by byte",
			func() *Merger {
				m := &Merger{}
				m.inputs = []io.Reader{
					strings.NewReader("abcd"),
					strings.NewReader("1234"),
					strings.NewReader("zxcv"),
				}
				return m
			},
			args{
				nil,
				1,
			},
			bytes.NewBuffer(make([]byte, 0, 12)),
			"a1zb2xc3cd4v",
			false,
		},
		{
			"Merge three files word by word",
			func() *Merger {
				m := &Merger{}
				m.inputs = []io.Reader{
					strings.NewReader("abcd"),
					strings.NewReader("1234"),
					strings.NewReader("zxcv"),
				}
				return m
			},
			args{
				nil,
				2,
			},
			bytes.NewBuffer(make([]byte, 0, 12)),
			"ab12zxcd34cv",
			false,
		},
		{
			"Merge three files dword by dword",
			func() *Merger {
				m := &Merger{}
				m.inputs = []io.Reader{
					strings.NewReader("abcddcba"),
					strings.NewReader("12344321"),
					strings.NewReader("zxcvvcxz"),
				}
				return m
			},
			args{
				nil,
				4,
			},
			bytes.NewBuffer(make([]byte, 0, 24)),
			"abcd1234zxcvdcba4321vcxz",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.prepare()
			m.output = tt.out
			if err := m.mergeSize(tt.args.ctx, tt.args.size); (err != nil) != tt.wantErr {
				t.Errorf("Merger.mergeSize() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.want, tt.out.String())
		})
	}
}

func TestMerger_bits(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		prepare func() *Merger
		args    args
		out     *bytes.Buffer
		want    []byte
		wantErr bool
	}{
		{
			"Merge two files bit by bit",
			func() *Merger {
				m := &Merger{}
				m.inputs = []io.Reader{
					bytes.NewReader([]byte{0b1111_0000}),
					bytes.NewReader([]byte{0b0000_1111}),
				}
				return m
			},
			args{nil},
			bytes.NewBuffer(make([]byte, 0, 4)),
			[]byte{0b0101_0101, 0b1010_1010},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := tt.prepare()
			m.output = tt.out
			if err := m.bits(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Merger.bits() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			require.Equal(t, tt.want, tt.out.Bytes())
		})
	}
}
