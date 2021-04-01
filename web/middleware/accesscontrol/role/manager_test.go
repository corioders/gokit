package role

import (
	"strconv"
	"testing"

	"github.com/corioders/gokit/errors"
)

func TestNewManager(t *testing.T) {
	for i := 0; i < 10; i++ {
		m1, err := NewManager("TestNewManager, roleManager, number of iterations " + strconv.Itoa(i))
		_ = m1
		if err != nil {
			t.Fatalf("Error while creating new roleManager, iteration: %v error: %v", i, err)
		}
	}

	m2, err := NewManager("TestNewManager, roleManager, name duplicate")
	_ = m2
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %v", err)
	}

	m2Duplicate, err := NewManager("TestNewManager, roleManager, name duplicate")
	_ = m2Duplicate
	if err != nil {
		if !errors.Is(err, ErrRoleManagerNonUnique) {
			t.Fatalf("Expected ErrRoleManagerNonUnique when roleManager name is a duplicate, but got: %v", err)
		}
	} else {
		t.Fatal("Expected error when creating new roleManager with duplicate name")
	}
}
