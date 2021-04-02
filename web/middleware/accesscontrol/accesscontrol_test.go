package accesscontrol

import (
	"io"
	"strconv"
	"testing"

	"github.com/corioders/gokit/errors"
	"github.com/corioders/gokit/math/rand"
)

var validAccesscontrolKey Key
var randReader *rand.Rand = rand.NewMath(0)

func init() {
	validAccesscontrolKey = make([]byte, 64)
	_, err := io.ReadFull(randReader, validAccesscontrolKey)
	if err != nil {
		panic(err)
	}
}

func TestNewAccesscontrol(t *testing.T) {
	for i := 0; i < 10; i++ {
		a1, err := New("TestNewAccesscontrol, accesscontrol, number of iterations "+strconv.Itoa(i), validAccesscontrolKey)
		_ = a1
		if err != nil {
			t.Fatalf("Expected no error while creating new accesscontrol, but on iteration: %v, got error: %v", i, err)
		}
	}

	a2, err := New("TestNewAccesscontrol, accesscontrol, name duplicate", validAccesscontrolKey)
	_ = a2
	if err != nil {
		t.Fatalf("Error while creating new accesscontrol, error: %v", err)
	}

	a2Duplicate, err := New("TestNewAccesscontrol, accesscontrol, name duplicate", validAccesscontrolKey)
	_ = a2Duplicate
	if err != nil {
		if !errors.Is(err, ErrNameNonUnique) {
			t.Fatalf("Expected ErrNameNonUnique when accesscontrol name is a duplicate, but got: %v", err)
		}
	} else {
		t.Fatal("Expected error when creating new accesscontrol with duplicate name")
	}

	invalidLengthAccesscontrolKey := make([]byte, 10)
	_, err = io.ReadFull(randReader, invalidLengthAccesscontrolKey)
	if err != nil {
		t.Fatalf("Error while reading into invalidAccesscontrolKey, error: %v", err)
	}

	a3, err := New("TestNewAccesscontrol, accesscontrol number 1", invalidLengthAccesscontrolKey)
	_ = a3
	if err != nil {
		if !errors.Is(err, ErrInvalidKeyLength) {
			t.Fatalf("Expected ErrInvalidKeyLength when accesscontrol key length is invalid, but got: %v", err)
		}
	} else {
		t.Fatal("Expected error when creating new accesscontrol with invalid key length")
	}
}
