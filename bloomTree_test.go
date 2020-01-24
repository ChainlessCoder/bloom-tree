package bloomtree

import (
	"fmt"
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
		{
			elements: [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}, {11}, {12}, {13},
				{14}, {15}, {16}},
			hashAt: [3]int{2, 3, 17},
		},
		{
			elements: [][]byte{{0}, {1}},
			hashAt:   [3]int{28, 29, 30},
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

func TestGetBloomFilter(t *testing.T) {
	var tests = []struct {
		elements [][]byte
	}{
		{
			elements: [][]byte{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}},
		},
		{
			elements: [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}, {11}, {12}, {13},
				{14}, {15}, {16}},
		},
		{
			elements: [][]byte{{0}, {1}},
		},
	}

	for _, test := range tests {
		dbf := generateDBF(test.elements...)
		tree, err := NewBloomTree(dbf)
		if err != nil {
			t.Fatal(err)
		}

		if !tree.GetBloomFilter().BitArray().Equal(dbf.BitArray()) {
			t.Fatal("bloom filter is not equal")
		}
	}
}

func TestBloomTreeExceedingK(t *testing.T) {
	dbf := DBF.NewDbf(200, 1e-100, []byte("secret seed"))
	_, err := NewBloomTree(dbf)
	if err.Error() != fmt.Errorf("parameter k of the bloom filter must be smaller than %d", maxK).Error() {
		t.Fatalf("expected error %v", fmt.Errorf("parameter k of the bloom filter must be smaller than %d", maxK))
	}
}

func generateDBF(elements ...[]byte) *DBF.DistBF {
	dbf := DBF.NewDbf(200, 0.2, []byte("secret seed"))
	for _, elem := range elements {
		dbf.Add(elem)
	}
	return dbf
}
