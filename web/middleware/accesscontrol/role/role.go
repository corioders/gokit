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

	role, err := manager.RoleFromName(rrj.Name)
	if err != nil {
		return errors.WithStack(err)
	}

	r.ri = role.ri
	return nil
}

func newRole(name, managerName string, permissions ...*Permission) *roleInternal {
	ri := &roleInternal{
		name:        name,
		managerName: managerName,

		rwmu:        sync.RWMutex{},
		permissions: make(map[string]*Permission),
	}

	for _, permission := range permissions {
		ri.permissions[permission.name] = permission
	}

	return ri
}

// NewRole creates new role, name must be unique across one RoleManager.
// NewRole return error if name is non unique.
func (rm *RoleManager) NewRole(name string, permissions ...*Permission) (*Role, error) {
	ri := newRole(name, rm.name, permissions...)
	role := &Role{
		ri: ri,
	}

	err := rm.insertRole(role)
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

// RoleFromName returns role associated with given name.
// RoleFromName return error if role with given name does not exists.
// Note that role which is marshaled as json is the same as role's name.
func (rm *RoleManager) RoleFromName(name string) (*Role, error) {
	role, ok := rm.getRole(name)
	if !ok {
		return nil, errors.WithMessage(ErrRoleNotExists, fmt.Sprintf(`role name "%v"`, name))
	}

	return role, nil
}

// Assign a permission to the role.
func (r *Role) Allow(p *Permission) {
	r.ri.rwmu.Lock()
	r.ri.permissions[p.name] = p
	r.ri.rwmu.Unlock()
}

// Revoke the specific permission.
func (r *Role) Disallow(p *Permission) {
	r.ri.rwmu.Lock()
	delete(r.ri.permissions, p.name)
	r.ri.rwmu.Unlock()
}

// Permit returns true if the role has specific permission.
func (r *Role) IsAllowedTo(p *Permission) bool {
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
