package sbt

import (
	"testing"
	"fmt"
	"bytes"
	"github.com/labbloom/DBF"
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

	// Scenario - Check behavior given indices for a missing element
	dbf = DBF.NewDbf(200, 0.1, []byte("seed"))
	allIndices = []int{0,1,2,956,957,958}
	indices = []int{1,955}
	dbf.SetIndices(allIndices)
	bT = NewBloomTree(dbf)
	proof, data, err = bT.generatePresenceProof(indices)
	// Check err
	if err == nil {
		t.Fatal(err)
	}
	// Check data
	d := []byte(fmt.Sprintf("3,%d", 0))
	res := bytes.Compare(d,data[0])
	if  res != 0 {
		t.Fatalf("expected %v, got %v", d, data[0])
	}
	// Check proof
	if proof != nil {
		t.Fatal("Multiproof was not able to get computed")
	}

}

func TestGenerateAbsenceProof(t *testing.T) {
	// Scenario - Check behavior given a missing value at the beggining of the tree
	dbf := DBF.NewDbf(200, 0.1, []byte("seed"))
	allIndices := []int{5, 10, 100}
	index := 2
	dbf.SetIndices(allIndices)
	bT := NewBloomTree(dbf)
	proof, data, err := bT.generateAbsenceProof(index)
	// Check err
	if err != nil {
		t.Fatal(err)
	}
	// Check data
	d := []byte(fmt.Sprintf("1,%d", allIndices[0]))
	res := bytes.Compare(d,data[0])
	if  res != 0 {
		t.Fatalf("expected %v, got %v", d, data[0])
	}
	
	// Check proof
	if proof == nil {
		t.Fatal("Multiproof was not able to get computed")
	}

	// Scenario - Check behavior given a missing value at the end of the tree
	dbf = DBF.NewDbf(200, 0.1, []byte("seed"))
	allIndices = []int{5, 10, 100}
	index = 200
	dbf.SetIndices(allIndices)
	bT = NewBloomTree(dbf)
	proof, data, err = bT.generateAbsenceProof(index)
	// Check err
	if err != nil {
		t.Fatal(err)
	}
	// Check data
	d = []byte(fmt.Sprintf("1,%d", 100))
	res = bytes.Compare(d,data[0])
	if  res != 0 {
		t.Fatalf("expected %v, got %v", d, data[0])
	}
	// Check proof
	if proof == nil {
		t.Fatal("Multiproof was not able to get computed")
	}

	// Scenario - Check behavior given a missing value somewhere in the middle of the tree
	dbf = DBF.NewDbf(200, 0.1, []byte("seed"))
	allIndices = []int{5, 10, 100}
	index = 8
	dbf.SetIndices(allIndices)
	bT = NewBloomTree(dbf)
	proof, data, err = bT.generateAbsenceProof(index)
	// Check err
	if err != nil {
		t.Fatal(err)
	}
	// Check data
	expectedIndices := []int{5,10}
	for ind, v := range data {
		d = []byte(fmt.Sprintf("1,%d", expectedIndices[ind]))
		res = bytes.Compare(d,v)
		if  res != 0 {
			t.Fatalf("expected %v, got %v", d, v)
		}
	}
	// Check proof
	if proof == nil {
		t.Fatal("Multiproof was not able to get computed")
	}
	
}


func TestSbtMultiProof(t *testing.T) {
	// Scenario - Check behavior given an element that is present in the DBF
	dbf := DBF.NewDbf(200, 0.2, []byte("seed"))
	element := []byte("something")
	dbf.Add(element)
	bT := NewBloomTree(dbf)
	proof, data, presence, err := bT.SbtMultiProof(element)
	// Check err
	if err != nil {
		t.Fatal(err)
	}
	// Check presence / absence
	if presence != true {
		t.Fatal(err)
	}
	// Check data
	expectedIndices, _ := dbf.Proof(element)
	for ind, v := range data {
		d := []byte(fmt.Sprintf("1,%d", expectedIndices[ind]))
		res := bytes.Compare(d,v)
		if  res != 0 {
			t.Fatalf("expected %v, got %v", d, v)
		}
	}
	// Check proof
	if proof == nil {
		t.Fatal("Multiproof was not able to get computed")
	}

	// Scenario - Check behavior given an element that is absent in the DBF
	dbf = DBF.NewDbf(200, 0.2, []byte("seed"))
	element = []byte("something")
	element1 := []byte("something else")
	dbf.Add(element)
	bT = NewBloomTree(dbf)
	proof, data, presence, err = bT.SbtMultiProof(element1)
	// Check err
	if err != nil {
		t.Fatal(err)
	}
	// Check presence / absence
	if presence == true {
		t.Fatal(err)
	}
	// Check data
	expectedIndices = []int{247, 484}
	for ind, v := range data {
		d := []byte(fmt.Sprintf("1,%d", expectedIndices[ind]))
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