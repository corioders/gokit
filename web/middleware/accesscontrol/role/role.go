package role

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/corioders/gokit/errors"
)

type roleInternal struct {
	name        string
	managerName string

	rwmu        sync.RWMutex
	permissions map[string]*Permission
}

type Role struct {
	ri *roleInternal
}

var (
	ErrNilPermission = errors.New("Nil permission pointer was passeed")
)

type roleJson struct {
	Name        string `json:"n"`
	ManagerName string `json:"mn"`
}

// MarshalJSON implements Marshaler interface.
func (r *Role) MarshalJSON() ([]byte, error) {
	rrj := roleJson{r.ri.name, r.ri.managerName}
	return json.Marshal(rrj)
}

// UnmarshalJSON implements Unmarshaler interface.
func (r *Role) UnmarshalJSON(data []byte) error {
	rrj := roleJson{}
	err := json.Unmarshal(data, &rrj)
	if err != nil {
		return errors.WithStack(err)
	}

	manager, err := managerFromName(rrj.ManagerName)
	if err != nil {
		return errors.WithStack(err)
	}

	role, ok := manager.getRole(rrj.Name)
	if !ok {
		return errors.WithMessage(ErrRoleNotExists, fmt.Sprintf(`role name: "%v"`, rrj.Name))
	}

	r.ri = role.ri
	return nil
}

func newRole(name, managerName string, permissions ...*Permission) (*roleInternal, error) {
	ri := &roleInternal{
		name:        name,
		managerName: managerName,

		rwmu:        sync.RWMutex{},
		permissions: make(map[string]*Permission),
	}

	for _, permission := range permissions {
		if permission == nil {
			return nil, errors.WithStack(ErrNilPermission)
		}
		ri.permissions[permission.name] = permission
	}

	return ri, nil
}

// NewRole creates new role, name must be unique across one RoleManager.
// NewRole return error if name is non unique.
func (rm *RoleManager) NewRole(name string, permissions ...*Permission) (*Role, error) {
	ri, err := newRole(name, rm.name, permissions...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	role := &Role{
		ri: ri,
	}

	err = rm.insertRole(role)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return role, nil
}

// MustNewRole is the same as NewRole, accept it panics is name is non unique.
func (rm *RoleManager) MustNewRole(name string, permissions ...*Permission) *Role {
	role, err := rm.NewRole(name, permissions...)
	if err != nil {
		panic(err)
	}

	return role
}

// Allow a permission to the role.
// Allow panics if any of permissions is nil.
func (r *Role) Allow(permissions ...*Permission) {
	r.ri.rwmu.Lock()
	for _, permission := range permissions {
		if permission == nil {
			r.ri.rwmu.Unlock()
			panic(errors.WithStack(ErrNilPermission))
		}
		r.ri.permissions[permission.name] = permission
	}
	r.ri.rwmu.Unlock()
}

// Disallow the specific permission.
// Disallow panics if any of permissions is nil.
func (r *Role) Disallow(permissions ...*Permission) {
	r.ri.rwmu.Lock()
	for _, permission := range permissions {
		if permission == nil {
			r.ri.rwmu.Unlock()
			panic(errors.WithStack(ErrNilPermission))
		}
		delete(r.ri.permissions, permission.name)
	}
	r.ri.rwmu.Unlock()
}

// IsAllowedTo returns true if the role has specific permission.
// IsAllowedTo panics if p is nil.
func (r *Role) IsAllowedTo(p *Permission) bool {
	if p == nil {
		panic(errors.WithStack(ErrNilPermission))
	}
	found := false
	r.ri.rwmu.RLock()
	for _, rp := range r.ri.permissions {
		if rp.matches(p) {
			found = true
			break
		}
	}
	r.ri.rwmu.RUnlock()
	return found
}

// func (r *Role) Matches(a *Role) bool {
// 	return r.ri.name == a.ri.name && r.ri.managerName == a.ri.managerName
// }
