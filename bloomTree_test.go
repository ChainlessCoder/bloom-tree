package bloomtree

import (
	"testing"

	"github.com/labbloom/DBF"
)

func TestNewBloomTree(t *testing.T) {
	var tests = []struct {
		elements [][]byte
		hashAt   [3]int
	}{
		{
			elements: [][]byte{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}},
			hashAt:   [3]int{0, 1, 16},
		},
	}

	for _, test := range tests {
		dbf := generateDBF(test.elements...)
		tree, err := NewBloomTree(dbf)
		if err != nil {
			t.Fatal(err)
		}

		if hashChild(tree.nodes[test.hashAt[0]], tree.nodes[test.hashAt[1]]) != tree.nodes[test.hashAt[2]] {
			t.Fatalf("h(%d, %d) != %d", test.hashAt[0], test.hashAt[1], test.hashAt[2])
		}
	}
}

func generateDBF(elements ...[]byte) *DBF.DistBF {
	dbf := DBF.NewDbf(200, 0.2, []byte("secret seed"))
	for _, elem := range elements {
		dbf.Add(elem)
	}
	return dbf
}
