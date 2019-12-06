package sbt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/willf/bitset"
)

func TestBit2Int(t *testing.T) {
	b := bitset.New(37)
	b.Set(0)
	b.Set(1)
	b.Set(8)
	b.Set(13)
	b.Set(18)
	b.Set(19)
	b.Set(20)
	b.Set(26)
	b.Set(28)
	b.Set(32)
	b.Set(33)
	b.Set(37)
	var tests = []struct {
		b             *bitset.BitSet
		expectedState [][2]int
	}{
		{b: b, expectedState: [][2]int{{2, 0}, {1, 8}, {1, 13}, {3, 18}, {1, 26}, {1, 28}, {2, 32}, {1, 37}}},
	}

	for i, test := range tests {
		state := bit2int(test.b)
		assert.Equal(t, test.expectedState, state, fmt.Sprintf("expected states to be the same at test %d", i))

	}
}
