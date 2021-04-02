package role

import "github.com/corioders/gokit/errors"

type Permission struct {
	name        string
	managerName string
}

func newPermission(name, managerName string) *Permission {
	return &Permission{
		name:        name,
		managerName: managerName,
	}
}

// NewPermission creates new permission, name must be unique across one RoleManager.
// NewPermission return error if name is non unique.
func (rm *RoleManager) NewPermission(name string) (*Permission, error) {
	permission := newPermission(name, rm.name)

	err := rm.insertPermission(permission)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return permission, nil
}

// MustNewPermission is the same as NewPermission, accept it panics is name is non unique.
func (rm *RoleManager) MustNewPermission(name string) *Permission {
	permission, err := rm.NewPermission(name)
	if err != nil {
		panic(err)
	}

	return permission
}

func (p *Permission) matches(a *Permission) bool {
	return p.name == a.name && p.managerName == a.managerName
}
