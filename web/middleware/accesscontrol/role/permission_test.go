package role

import (
	"strconv"
	"testing"

	"github.com/corioders/gokit/errors"
)

func TestNewPermission(t *testing.T) {
	roleManager, err := NewManager("TestNewPermission, roleManager, number 1")
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %e", err)
	}

	for i := 0; i < 10; i++ {
		p1, err := roleManager.NewPermission("TestNewPermission, permission, number of iterations " + strconv.Itoa(i))
		_ = p1
		if err != nil {
			t.Fatalf("Error while creating new permission, iteration: %v error: %v", i, err)
		}
	}

	p2, err := roleManager.NewPermission("TestNewPermission, permission, name duplicate")
	_ = p2
	if err != nil {
		t.Fatalf("Error while creating new permission, error: %v", err)
	}

	p2Duplicate, err := roleManager.NewPermission("TestNewPermission, permission, name duplicate")
	_ = p2Duplicate
	if err != nil {
		if !errors.Is(err, ErrPermissionNameNonUnique) {
			t.Fatalf("Expected ErrPermissionNameNonUnique when permission name is a duplicate, but got: %v", err)
		}
	} else {
		t.Fatal("Expected error when creating new permission with duplicate name")
	}
}
