package role

import (
	"fmt"
	"sync"

	"github.com/corioders/gokit/errors"
)

var roleManagers = &sync.Map{}

type RoleManager struct {
	name string

	roles       *sync.Map
	permissions *sync.Map
}

var (
	ErrRoleManagerNonUnique    = errors.New("Role manager name must be unique")
	ErrRoleNameNonUnique       = errors.New("Role name must be unique")
	ErrPermissionNameNonUnique = errors.New("Permission name must be unique")

	ErrRoleNotExists = errors.New("Role does not exist")

	// ErrRoleManagerNotExists occurs only when in json manager name is not valid
	ErrRoleManagerNotExists = errors.New("Role manager does not exist")
)

// NewManager creates new permission, name must be alway unique.
// NewManager return error if name is non unique.
func NewManager(name string) (*RoleManager, error) {
	_, ok := roleManagers.Load(name)
	if ok {
		return nil, errors.WithMessage(ErrRoleManagerNonUnique, fmt.Sprintf(`name "%v" is not unique`, name))
	}

	roleManager := &RoleManager{
		name: name,

		roles:       &sync.Map{},
		permissions: &sync.Map{},
	}
	roleManagers.Store(name, roleManager)

	return roleManager, nil
}

// MustNewManager is the same as NewManager, accept it panics is name is non unique.
func MustNewManager(name string) *RoleManager {
	roleManager, err := NewManager(name)
	if err != nil {
		panic(err)
	}

	return roleManager
}

func managerFromName(name string) (*RoleManager, error) {
	roleManager, ok := roleManagers.Load(name)
	if !ok {
		return nil, ErrRoleManagerNotExists
	}

	return roleManager.(*RoleManager), nil
}

func (rm *RoleManager) insertRole(r *Role) error {
	_, ok := rm.roles.Load(r.ri.name)
	if ok {
		return errors.WithMessage(ErrRoleNameNonUnique, fmt.Sprintf(`name "%v" is not unique`, r.ri.name))
	}

	rm.roles.Store(r.ri.name, r)
	return nil
}

func (rm *RoleManager) getRole(name string) (*Role, bool) {
	role, ok := rm.roles.Load(name)
	if !ok {
		return nil, false
	}

	return role.(*Role), true
}

func (rm *RoleManager) insertPermission(p *Permission) error {
	_, ok := rm.permissions.Load(p.name)
	if ok {
		return errors.WithMessage(ErrPermissionNameNonUnique, fmt.Sprintf(`name "%v" is not unique`, p.name))
	}

	rm.permissions.Store(p.name, p)
	return nil
}
