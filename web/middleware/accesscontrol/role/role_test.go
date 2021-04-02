package role

import (
	"strconv"
	"testing"

	"github.com/corioders/gokit/errors"
)

func TestNewRole(t *testing.T) {
	roleManager, err := NewManager("TestNewRole, roleManager, number 1")
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %v", err)
	}

	for i := 0; i < 10; i++ {
		r1, err := roleManager.NewRole("TestNewRole, role, number of iterations " + strconv.Itoa(i))
		_ = r1
		if err != nil {
			t.Fatalf("Expected no error while creating new role, but on iteration: %v, got error: %v", i, err)
		}
	}

	r2, err := roleManager.NewRole("TestNewRole, role, name duplicate")
	_ = r2
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	r2Duplicate, err := roleManager.NewRole("TestNewRole, role, name duplicate")
	_ = r2Duplicate
	if err != nil {
		if !errors.Is(err, ErrRoleNameNonUnique) {
			t.Fatalf("Expected ErrRoleNameNonUnique when role name is a duplicate, but got: %v", err)
		}
	} else {
		t.Fatal("Expected error when creating new role with duplicate name")
	}

	r3, err := roleManager.NewRole("TestNewRole, role, number 3", nil)
	_ = r3
	if err != nil {
		if !errors.Is(err, ErrNilPermission) {
			t.Fatalf("Expected ErrNilPermission when permission passed to NewRole is nil, but got: %v", err)
		}
	} else {
		t.Fatal("Expected error when createing new role with nil permission")
	}
}

func TestIsAllowedTo(t *testing.T) {
	roleManager, err := NewManager("TestIsAllowedTo, roleManager, number 1")
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %v", err)
	}

	p1, err := roleManager.NewPermission("TestIsAllowedTo, permission, number 1")
	if err != nil {
		t.Fatalf("Error while creating new permission error: %v", err)
	}

	p2, err := roleManager.NewPermission("TestIsAllowedTo, permission, number 2")
	if err != nil {
		t.Fatalf("Error while creating new permission error: %v", err)
	}

	r1, err := roleManager.NewRole("TestIsAllowedTo, role, number 1", p1)
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	if !r1.IsAllowedTo(p1) {
		t.Fatal(`Expected "r1" to be allowed to "p1", because when "r1" is created "p1" is passed`)
	}

	if r1.IsAllowedTo(p2) {
		t.Fatal(`Expected "r1" to not be allowed to "p2", because when "r1" is created "p2" is not passed`)
	}

	r2, err := roleManager.NewRole("TestIsAllowedTo, role, number 2", p1, p2)
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	if !r2.IsAllowedTo(p1) {
		t.Fatal(`Expected "r2" to be allowed to "p1", because when "r2" is created "p1" is passed`)
	}

	if !r2.IsAllowedTo(p2) {
		t.Fatal(`Expected "r2" to be allowed to "p2", because when "r2" is created "p2" is passed`)
	}

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("Expected panic when nil permission passed to IsAllowedTo")
			}

			err, ok := r.(error)
			if !ok {
				t.Fatalf("Expected recover value to be of type error, but got: %v", r)
			}

			if !errors.Is(err, ErrNilPermission) {
				t.Fatal("Expected ErrNilPermission when nil permission passed to IsAllowedTo")
			}
		}()

		r1.IsAllowedTo(nil)
	}()
}

func TestAllow(t *testing.T) {
	roleManager, err := NewManager("TestAllow, roleManager, number 1")
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %v", err)
	}

	p1, err := roleManager.NewPermission("TestAllow, permission, number 1")
	if err != nil {
		t.Fatalf("Error while creating new permission error: %v", err)
	}

	p2, err := roleManager.NewPermission("TestAllow, permission, number 2")
	if err != nil {
		t.Fatalf("Error while creating new permission error: %v", err)
	}

	r1, err := roleManager.NewRole("TestAllow, role, number 1")
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	r1.Allow(p1)
	if !r1.IsAllowedTo(p1) {
		t.Fatal(`Expected "r1" to be allowed to "p1", because when r1.Allow is called "p1" is passed`)
	}

	if r1.IsAllowedTo(p2) {
		t.Fatal(`Expected "r1" to not be allowed to "p2", because when "r1" is created "p2" is not passed`)
	}

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("Expected panic when nil permission passed to Allow")
			}

			err, ok := r.(error)
			if !ok {
				t.Fatalf("Expected recover value to be of type error, but got: %v", r)
			}

			if !errors.Is(err, ErrNilPermission) {
				t.Fatal("Expected ErrNilPermission when nil permission passed to Allow")
			}
		}()

		r1.Allow(nil)
	}()

	r2, err := roleManager.NewRole("TestAllow, role, number 2")
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	r2.Allow(p1)
	if !r2.IsAllowedTo(p1) {
		t.Fatal(`Expected "r2" to be allowed to "p1", because when r2.Allow is called "p1" is passed`)
	}

	r2.Allow(p2)
	if !r2.IsAllowedTo(p2) {
		t.Fatal(`Expected "r2" to be allowed to "p2", because when r2.Allow is called "p2" is passed`)
	}
}

func TestDisallow(t *testing.T) {
	roleManager, err := NewManager("TestDisallow, roleManager, number 1")
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %v", err)
	}

	p1, err := roleManager.NewPermission("TestDisallow, permission, number 1")
	if err != nil {
		t.Fatalf("Error while creating new permission error: %v", err)
	}

	p2, err := roleManager.NewPermission("TestDisallow, permission, number 2")
	if err != nil {
		t.Fatalf("Error while creating new permission error: %v", err)
	}

	r1, err := roleManager.NewRole("TestDisallow, role, number 1", p1, p2)
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	r1.Disallow(p2)
	if !r1.IsAllowedTo(p1) {
		t.Fatal(`Expected "r1" to be allowed to "p1", because when "r1" is created "p1" is passed`)
	}

	if r1.IsAllowedTo(p2) {
		t.Fatal(`Expected "r1" to not be allowed to "p2", because "r1.Disallow" is called with "p2" passed`)
	}

	func() {
		defer func() {
			r := recover()
			if r == nil {
				t.Fatal("Expected panic when nil permission passed to Disallow")
			}

			err, ok := r.(error)
			if !ok {
				t.Fatalf("Expected recover value to be of type error, but got: %v", r)
			}

			if !errors.Is(err, ErrNilPermission) {
				t.Fatal("Expected ErrNilPermission when nil permission passed to Disallow")
			}
		}()

		r1.Disallow(nil)
	}()

	r2, err := roleManager.NewRole("TestDisallow, role, number 2", p1, p2)
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	r2.Disallow(p1)
	if r2.IsAllowedTo(p1) {
		t.Fatal(`Expected "r2" to not be allowed to "p1", because "r2.Disallow" is called with "p1" passed`)
	}

	r2.Disallow(p2)
	if r2.IsAllowedTo(p2) {
		t.Fatal(`Expected "r2" to not be allowed to "p2", because "r2.Disallow" is called with "p2" passed`)
	}
}
