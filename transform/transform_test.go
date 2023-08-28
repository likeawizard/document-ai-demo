package transform

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMoneyParser(t *testing.T) {

	delta := float64(0.005)

	type testCase struct {
		name      string
		stringVal string
		want      float64
	}

	tcs := []testCase{
		{
			name:      "Empty",
			stringVal: "",
			want:      0,
		},
		{
			name:      "Standard decimal",
			stringVal: "12.34",
			want:      12.34,
		},
		{
			name:      "Comma decimal",
			stringVal: "12,34",
			want:      12.34,
		},
		{
			name:      "Standard decimal, comma separator",
			stringVal: "12,122.34",
			want:      12122.34,
		},
		{
			name:      "Comma decimal, point separator",
			stringVal: "12.122,34",
			want:      12122.34,
		},
		{
			name:      "Non-numeric",
			stringVal: "About Tree Fiddy and a quarter",
			want:      0,
		},
	}

	for _, tc := range tcs {
		assert.InDelta(t, moneyParser(tc.stringVal), tc.want, delta, tc.name)
	}
}
