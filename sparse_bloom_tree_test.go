package sbt

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/labbloom/DBF"
)

func TestBit2Int(t *testing.T) {
	// Scenario - Check behavior given simple DBF (every positive integer is 1)
	dbf := DBF.NewDbf(100, 0.1, []byte("seed"))
	indices := []int{4,7,11}
	dbf.SetIndices(indices)
	b2i := bit2int(dbf.BitArray())
	expected := [][2]int{{1,4}, {1,7}, {1,11}}
	assert.Equal(t, expected, b2i, fmt.Sprintf("The 2D integer array is wrong"))

	// Scenario - Check behavior given more complex DBF (some positive integers can be greater than 1)
	dbf = DBF.NewDbf(100, 0.1, []byte("seed"))
	indices = []int{4,5,7,10,11,12}
	dbf.SetIndices(indices)
	b2i = bit2int(dbf.BitArray())
	expected = [][2]int{{2,4}, {1,7}, {3,10}}
	assert.Equal(t, expected, b2i, fmt.Sprintf("The 2D integer array is wrong"))
}
