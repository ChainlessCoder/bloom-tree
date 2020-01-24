package bloomtree

import (
	"crypto/sha512"
	"testing"
)

func TestHashLeaf(t *testing.T) {
	var tests = []struct {
		element uint64
		index   uint64
		output  [sha512.Size256]byte
	}{
		{
			element: 0,
			index:   1,
			output: [sha512.Size256]byte{173, 22, 23, 123, 203, 81, 76, 177, 198, 139, 98, 52, 238,
				132, 239, 103, 154, 150, 112, 4, 3, 5, 188, 227, 162, 24, 70, 210, 212, 90, 76, 51},
		},
	}

	for _, test := range tests {
		output := hashLeaf(test.element, test.index)
		if output != test.output {
			t.Fatalf("test failed at hashing element %d and index %d", test.element, test.index)
		}
	}
}

func TestHashChild(t *testing.T) {
	var tests = []struct {
		element1 [sha512.Size256]byte
		element2 [sha512.Size256]byte
		output   [sha512.Size256]byte
	}{
		{
			element1: hashLeaf(0, 1),
			element2: hashLeaf(1, 2),
			output: [sha512.Size256]byte{87, 189, 113, 25, 180, 188, 126, 186, 162, 0, 130, 60, 78, 222,
				215, 35, 16, 89, 51, 71, 195, 215, 115, 228, 110, 158, 24, 143, 163, 78, 246, 77},
		},
		{
			element1: hashLeaf(10, 11),
			element2: hashLeaf(11, 12),
			output: [sha512.Size256]byte{117, 149, 100, 207, 3, 124, 127, 199, 244, 22, 35, 97, 106, 31, 148,
				156, 177, 46, 21, 184, 174, 33, 110, 160, 66, 213, 109, 177, 2, 70, 48, 136},
		},
	}

	for _, test := range tests {
		output := hashChild(test.element1, test.element2)
		if output != test.output {
			t.Fatal("test failed at hashing child")
		}
	}
}
