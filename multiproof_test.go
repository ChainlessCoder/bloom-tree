package sbt

import (
	"testing"
	"fmt"
	"bytes"
	//"github.com/labbloom/go-merkletree"
	"github.com/labbloom/DBF"
	//"github.com/willf/bitset"
)


func TestGeneratePresenceProof(t *testing.T) {
	
	// Scenario - Check behavior given simple DBF (every positive integer is 1)
	dbf := DBF.NewDbf(200, 0.1, []byte("seed"))
	indices := []int{4,7,11}
	dbf.SetIndices(indices)
	bT := NewBloomTree(dbf)
	proof, data, err := bT.generatePresenceProof(indices)
	// Check err
	if err != nil {
		t.Fatal(err)
	}
	// Check data
	for ind, v := range data {
		d := []byte(fmt.Sprintf("1,%d", indices[ind]))
		res := bytes.Compare(d,v)
		if  res != 0 {
			t.Fatalf("expected %v, got %v", d, v)
		}
	}
	// Check proof
	if proof == nil {
		t.Fatal("Multiproof was not able to get computed")
	}

	// Scenario - Check behavior given more complex DBF (some positive integers are 2)
	dbf = DBF.NewDbf(200, 0.1, []byte("seed"))
	allIndices := []int{4,5,7,8,10,11}
	indices = []int{4,7,11}
	dbf.SetIndices(allIndices)
	bT = NewBloomTree(dbf)
	proof, data, err = bT.generatePresenceProof(indices)
	// Check err
	if err != nil {
		t.Fatal(err)
	}
	// Check data
	expectedIndices := []int{4, 7, 10}
	for ind, v := range data {
		d := []byte(fmt.Sprintf("2,%d", expectedIndices[ind]))
		res := bytes.Compare(d,v)
		if  res != 0 {
			t.Fatalf("expected %v, got %v", d, v)
		}
	}
	// Check proof
	if proof == nil {
		t.Fatal("Multiproof was not able to get computed")
	}

	// Scenario - Check behavior given consecutive ones at the beginning and end of DBF
	dbf = DBF.NewDbf(200, 0.1, []byte("seed"))
	allIndices = []int{0,1,2,956,957,958}
	indices = []int{1,958}
	dbf.SetIndices(allIndices)
	bT = NewBloomTree(dbf)
	proof, data, err = bT.generatePresenceProof(indices)
	// Check err
	if err != nil {
		t.Fatal(err)
	}
	// Check data
	expectedIndices = []int{0,956}
	for ind, v := range data {
		d := []byte(fmt.Sprintf("3,%d", expectedIndices[ind]))
		res := bytes.Compare(d,v)
		if  res != 0 {
			t.Fatalf("expected %v, got %v", d, v)
		}
	}
	// Check proof
	if proof == nil {
		t.Fatal("Multiproof was not able to get computed")
	}

}

