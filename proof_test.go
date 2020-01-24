package bloomtree

import (
	"errors"
	"testing"
)

func TestPresenceProofPresentElement(t *testing.T) {
	var tests = []struct {
		element  []byte
		elements [][]byte
	}{
		{
			element:  []byte{1},
			elements: [][]byte{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}},
		},
		{
			element: []byte{1},
			elements: [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}, {11}, {12}, {13},
				{14}, {15}, {16}},
		},
		{
			element:  []byte{1},
			elements: [][]byte{{0}, {1}},
		},
	}

	for _, test := range tests {
		seed := "secret seed"
		dbf := generateDBF(seed, test.elements...)
		tree, err := NewBloomTree(dbf)
		if err != nil {
			t.Fatal(err)
		}

		multiproof, err := tree.GenerateCompactMultiProof(test.element)
		if err != nil {
			t.Fatal(err)
		}

		if CheckProofType(multiproof.proofType) != true {
			t.Fatal("proof type is not presence")
		}

		present, err := tree.VerifyCompactMultiProof(test.element, []byte(seed), multiproof, tree.Root())
		if err != nil {
			t.Fatal(err)
		} else if !present {
			t.Fatal("expected element to be present, but is absent")
		}
	}
}

func TestPresenceProofAbsentElement(t *testing.T) {
	var tests = []struct {
		element      []byte
		elementProof []byte
		elements     [][]byte
	}{
		{
			element:      []byte{9},
			elementProof: []byte{8},
			elements:     [][]byte{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}},
		},
		{
			element:      []byte{17},
			elementProof: []byte{16},
			elements: [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}, {11}, {12}, {13},
				{14}, {15}, {16}},
		},
		{
			element:      []byte{2},
			elementProof: []byte{0},
			elements:     [][]byte{{0}, {1}},
		},
	}

	for _, test := range tests {
		seed := "secret seed"
		dbf := generateDBF(seed, test.elements...)
		tree, err := NewBloomTree(dbf)
		if err != nil {
			t.Fatal(err)
		}

		multiproof, err := tree.GenerateCompactMultiProof(test.elementProof)
		if err != nil {
			t.Fatal(err)
		}

		if CheckProofType(multiproof.proofType) != true {
			t.Fatal("proof type is not present")
		}

		_, err = tree.VerifyCompactMultiProof(test.element, []byte(seed), multiproof, tree.Root())
		if err == nil {
			t.Fatalf("expected error: %v", errors.New("the element is not inside the provided chunks for a presence proof"))
		} else if err.Error() != errors.New("the element is not inside the provided chunks for a presence proof").Error() {
			t.Fatalf("expected error %v, but got %v", errors.New("the element is not inside the provided chunks for a presence proof"), err)
		}
	}
}

func TestAbsentProofAbsentElement(t *testing.T) {
	var tests = []struct {
		element  []byte
		elements [][]byte
	}{
		{
			element:  []byte{9},
			elements: [][]byte{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}},
		},
		{
			element: []byte{17},
			elements: [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}, {11}, {12}, {13},
				{14}, {15}, {16}},
		},
		{
			element:  []byte{2},
			elements: [][]byte{{0}, {1}},
		},
	}

	for _, test := range tests {
		seed := "secret seed"
		dbf := generateDBF(seed, test.elements...)
		tree, err := NewBloomTree(dbf)
		if err != nil {
			t.Fatal(err)
		}

		multiproof, err := tree.GenerateCompactMultiProof(test.element)
		if err != nil {
			t.Fatal(err)
		}

		if CheckProofType(multiproof.proofType) != false {
			t.Fatal("proof type is not absent")
		}

		absent, err := tree.VerifyCompactMultiProof(test.element, []byte(seed), multiproof, tree.Root())
		if err != nil {
			t.Fatal(err)
		} else if !absent {
			t.Fatal("expected element to be absent, but is present")
		}
	}
}

func TestAbsenceProofPresentElement(t *testing.T) {
	var tests = []struct {
		element      []byte
		elementProof []byte
		elements     [][]byte
	}{
		{
			element:      []byte{8},
			elementProof: []byte{9},
			elements:     [][]byte{{1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}},
		},
		{
			element:      []byte{16},
			elementProof: []byte{17},
			elements: [][]byte{{0}, {1}, {2}, {3}, {4}, {5}, {6}, {7}, {8}, {9}, {10}, {11}, {12}, {13},
				{14}, {15}, {16}},
		},
		{
			element:      []byte{0},
			elementProof: []byte{2},
			elements:     [][]byte{{0}, {1}},
		},
	}

	for _, test := range tests {
		seed := "secret seed"
		dbf := generateDBF(seed, test.elements...)
		tree, err := NewBloomTree(dbf)
		if err != nil {
			t.Fatal(err)
		}

		multiproof, err := tree.GenerateCompactMultiProof(test.elementProof)
		if err != nil {
			t.Fatal(err)
		}

		if CheckProofType(multiproof.proofType) != false {
			t.Fatal("proof type is not absent")
		}

		_, err = tree.VerifyCompactMultiProof(test.element, []byte(seed), multiproof, tree.Root())
		if err != nil {
			t.Fatal(err)
		}
	}
}
