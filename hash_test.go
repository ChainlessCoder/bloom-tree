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
			output: [sha512.Size256]byte{79, 42, 198, 69, 197, 164, 159, 73, 97, 201, 36, 127, 235, 9, 221,
				214, 118, 111, 196, 191, 196, 127, 3, 212, 108, 204, 175, 4, 99, 143, 60, 51},
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
			output: [sha512.Size256]byte{202, 116, 135, 95, 85, 135, 228, 38, 153, 127, 237, 234, 194, 152, 113, 112,
				70, 226, 250, 42, 106, 63, 161, 138, 85, 110, 34, 240, 186, 151, 198, 108},
		},
		{
			element1: hashLeaf(10, 11),
			element2: hashLeaf(11, 12),
			output: [sha512.Size256]byte{105, 250, 104, 250, 231, 6, 222, 161, 109, 46, 208, 106, 94, 20, 246, 171, 169,
				116, 12, 124, 101, 111, 87, 91, 173, 114, 53, 89, 156, 86, 109, 190},
		},
	}

	for _, test := range tests {
		output := hashChild(test.element1, test.element2)
		if output != test.output {
			t.Fatal("test failed at hashing child")
		}
	}
}
