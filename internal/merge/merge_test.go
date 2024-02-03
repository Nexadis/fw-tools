package merge

import (
	"bytes"
	"context"
	"io"
	"os"
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
		{
			"Simple new",
			args{config.Merge{}},
			New(config.Merge{}),
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

func TestMerger_Run(t *testing.T) {
	type args struct {
		ctx    context.Context
		input1 string
		input2 string
	}
	tests := []struct {
		name    string
		prepare func() *Merger
		args    args
		out     []byte
		wantErr bool
	}{
		{
			"Split by byte 2 files",
			func() *Merger {
				cfg := config.Merge{
					ByByte: true,
				}

				m := New(cfg)
				return m
			},
			args{
				nil,
				"abcdefgh",
				"12345678",
			},
			[]byte("a1b2c3d4e5f6g7h8"),
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m *Merger
			if tt.prepare != nil {
				m = tt.prepare()
			}
			f1, _ := os.CreateTemp("", "*.bin")
			f1.Write([]byte(tt.args.input1))
			f1.Close()
			f2, _ := os.CreateTemp("", "*.bin")
			f2.Write([]byte(tt.args.input2))
			f2.Close()
			inputs := []string{f1.Name(), f2.Name()}
			fo, _ := os.CreateTemp("", "*.bin")
			fo.Close()
			output := fo.Name()
			err, cls := m.Open(inputs, output)
			require.NoError(t, err)
			if err := m.Run(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Merger.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
			cls()
			data, err := os.ReadFile(output)
			require.NoError(t, err)
			require.Equal(t, tt.out, data)
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
					bytes.NewReader([]byte{0b1111_1010}),
					bytes.NewReader([]byte{0b0000_1111}),
				}
				return m
			},
			args{nil},
			bytes.NewBuffer(make([]byte, 0, 2)),
			[]byte{0b1110_1110, 0b0101_0101},
			false,
		},
		{
			"Merge two files bit by bit",
			func() *Merger {
				m := &Merger{}
				m.inputs = []io.Reader{
					bytes.NewReader([]byte{0b1111_1010}),
					bytes.NewReader([]byte{0b0000_1111}),
					bytes.NewReader([]byte{0b0000_1111}),
				}
				return m
			},
			args{nil},
			bytes.NewBuffer(make([]byte, 0, 3)),
			[]byte{0b1011_1110, 0b1001_1111, 0b0010_0100},
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
