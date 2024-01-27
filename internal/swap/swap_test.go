package swap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwapBits(t *testing.T) {
	tests := []struct {
		name string
		arg  uint8
		want uint8
	}{
		{
			name: "Simple swap bits",
			arg:  0b10111001,
			want: 0b10011101,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, swapBits(test.arg))
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Equal(t, test.want, swapHalf(test.arg))
		})
	}
}
