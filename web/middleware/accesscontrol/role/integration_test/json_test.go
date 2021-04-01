package integration_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/corioders/gokit/web/middleware/accesscontrol/role"
)

func TestDirectJson(t *testing.T) {
	roleManager, err := role.NewManager("TestDirectJson, roleManager, number 1")
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %e", err)
	}

	normalPermission, err := roleManager.NewPermission("TestDirectJson, permission, normal premissions")
	if err != nil {
		t.Fatalf("Error while creating new permission, error: %v", err)
	}

	notNormalPermission, err := roleManager.NewPermission("TestDirectJson, permission, not normal premissions")
	if err != nil {
		t.Fatalf("Error while creating new permission, error: %v", err)
	}

	normalRole, err := roleManager.NewRole("TestDirectJson, role, normal role", normalPermission)
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	normalRoleJson, err := json.Marshal(normalRole)
	if err != nil {
		t.Fatalf(`Error while marshaling "normalRole" into "normalRoleJson", error: %v`, err)
	}

	normalRole2 := &role.Role{}
	err = json.Unmarshal(normalRoleJson, normalRole2)
	if err != nil {
		t.Fatalf(`Error while unmarshaling "normalRoleJson" into "normalRole2", error: %v`, err)
	}

	if !normalRole2.IsAllowedTo(normalPermission) {
		t.Fatal(`Expected "normalRole2" to be allowed to "normalPermission", because when "normalRole" is created "normalPermission" is passed, and normalRole2 is just normalRole marshaled and unmarshaled`)
	}

	if normalRole2.IsAllowedTo(notNormalPermission) {
		t.Fatal(`Expected "normalRole2" to not be allowed to "notNormalPermission", because when "normalRole" is created "notNormalPermission" is not passed, and normalRole2 is just normalRole marshaled and unmarshaled`)
	}
}

func TestEmbeddedJson(t *testing.T) {
	type testEmbeddedJson struct {
		Role *role.Role `json:"role"`
	}

	roleManager, err := role.NewManager("TestEmbeddedJson, roleManager, number 1")
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %e", err)
	}

	normalPermission, err := roleManager.NewPermission("TestEmbeddedJson, permission, normal premissions")
	if err != nil {
		t.Fatalf("Error while creating new permission, error: %v", err)
	}

	notNormalPermission, err := roleManager.NewPermission("TestEmbeddedJson, permission, not normal premissions")
	if err != nil {
		t.Fatalf("Error while creating new permission, error: %v", err)
	}

	normalRole, err := roleManager.NewRole("TestEmbeddedJson, role, normal role", normalPermission)
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	embedded := testEmbeddedJson{Role: normalRole}
	embeddedJson, err := json.Marshal(embedded)
	if err != nil {
		t.Fatalf(`Error while marshaling "embedded" into "embeddedJson", error: %v`, err)
	}

	embedded2 := testEmbeddedJson{}
	err = json.Unmarshal(embeddedJson, &embedded2)
	if err != nil {
		t.Fatalf(`Error while unmarshaling "embeddedJson" into "embedded2", error: %v`, err)
	}

	normalRole2 := embedded2.Role
	if !normalRole2.IsAllowedTo(normalPermission) {
		t.Fatal(`Expected "normalRole2" to be allowed to "normalPermission", because when "normalRole" is created "normalPermission" is passed, and normalRole2 is just normalRole marshaled and unmarshaled`)
	}

	if normalRole2.IsAllowedTo(notNormalPermission) {
		t.Fatal(`Expected "normalRole2" to not be allowed to "notNormalPermission", because when "normalRole" is created "notNormalPermission" is not passed, and normalRole2 is just normalRole marshaled and unmarshaled`)
	}
}

func TestBrokenJson(t *testing.T) {
	wrongManagerName := []byte(`{"n":"TestBrokenJson, role, not existing role","mn":"TestBrokenJson, roleManager, not existing roleManager"}`)
	err := json.Unmarshal(wrongManagerName, &role.Role{})
	if err == nil {
		t.Fatal("Expected error when unmarshaling broken role json")
	}

	if !errors.Is(err, role.ErrRoleManagerNotExists) {
		t.Fatal("Expected ErrRoleManagerNotExists when unmarshaling broken role json, and roleManager name is invalid")
	}

	roleManager, err := role.NewManager("TestBrokenJson, roleManager, number 1")
	_ = roleManager
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %e", err)
	}

	wrongRoleName := []byte(`{"n":"TestBrokenJson, role, not existing role","mn":"TestBrokenJson, roleManager, number 1"}`)
	err = json.Unmarshal(wrongRoleName, &role.Role{})
	if err == nil {
		t.Fatal("Expected error when unmarshaling broken role json")
	}

	if !errors.Is(err, role.ErrRoleNotExists) {
		t.Fatal("Expected ErrRoleNotExists when unmarshaling broken role json, and roleManager name is valid and role name is invalid")
	}
}
