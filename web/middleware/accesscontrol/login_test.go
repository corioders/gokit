package accesscontrol

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corioders/gokit/errors"
	"github.com/corioders/gokit/web/middleware/accesscontrol/role"
)

func TestNewLogin(t *testing.T) {
	accesscontrol, err := New("TestNewLogin, accesscontrol, number 1", validAccesscontrolKey)
	if err != nil {
		t.Fatalf("Error while creating new accesscontrol, error: %v", err)
	}

	t.Run("ValidLah", func(t *testing.T) {
		lah := func(ctx context.Context, r *http.Request) (claims interface{}, roleGranted *role.Role, shouldLogin bool, err error) {
			return nil, nil, true, nil
		}

		_, err = accesscontrol.NewLogin(lah)
		if err != nil {
			t.Fatalf("Expected no error while creating new login handler, but got error: %v", err)
		}
	})

	t.Run("NilLah", func(t *testing.T) {
		_, err = accesscontrol.NewLogin(nil)
		if err != nil {
			if !errors.Is(err, ErrNilLoginAcceptHandler) {
				t.Fatalf("Expected ErrNilLoginAcceptHandler while creating new login handler with nil LoginAcceptHandler")
			}
		} else {
			t.Fatalf("Expected error while creating new login handler with nil LoginAcceptHandler")
		}
	})
}

func TestLoginHandler(t *testing.T) {
	t.Run("ValidLoginAcceptHandler", func(t *testing.T) {
		accesscontrol, err := New("TestLoginHandler, ValidLoginAcceptHandler, accesscontrol, number 1", validAccesscontrolKey)
		if err != nil {
			t.Fatalf("Error while creating new accesscontrol, error: %v", err)
		}

		roleManager, err := role.NewManager("TestLoginHandler, ValidLoginAcceptHandler, roleManager, number 1")
		if err != nil {
			t.Fatalf("Error while creating new roleManager, error: %v", err)
		}

		normalPermission, err := roleManager.NewPermission("TestLoginHandler, ValidLoginAcceptHandler, perimssion, number 1")
		if err != nil {
			t.Fatalf("Error while creating new permission, error: %v", err)
		}

		normalRole, err := roleManager.NewRole("TestLoginHandler, ValidLoginAcceptHandler, role, number 1", normalPermission)
		if err != nil {
			t.Fatalf("Error while creating new role, error: %v", err)
		}

		lah := func(ctx context.Context, r *http.Request) (claims interface{}, roleGranted *role.Role, shouldLogin bool, err error) {
			return nil, normalRole, true, nil
		}

		loginHandler, err := accesscontrol.NewLogin(lah)
		if err != nil {
			t.Fatalf("Error while creating new login handler, error: %v", err)
		}

		ctx := context.Background()
		r := httptest.NewRequest("GET", "http://foo.com", nil)
		rw := httptest.NewRecorder()

		err = loginHandler(ctx, rw, r)
		if err != nil {
			t.Fatalf("Expected no error while executing valid loginHandler, but got error: %v", err)
		}

		if rw.Code != http.StatusOK {
			t.Fatal("Expected status code after valid loginHandler to be http.StatusOK")
		}

		cookies := rw.Result().Cookies()
		if len(cookies) != 1 {
			t.Fatal("Expected cookie count to be 1 after successful execution of loginHandler")
		}

		tokenCookie := cookies[0]
		if tokenCookie.Name != accesscontrol.tokenCookieName {
			t.Fatal("Expected cookie name to be accesscontrol.tokenCookieName")
		}

		token := tokenCookie.Value
		_ = token
	})

	t.Run("InvalidLoginAcceptHandler", func(t *testing.T) {
		accesscontrol, err := New("TestLoginHandler, InvalidLoginAcceptHandler, accesscontrol, number 1", validAccesscontrolKey)
		if err != nil {
			t.Fatalf("Error while creating new accesscontrol, error: %v", err)
		}

		lah := func(ctx context.Context, r *http.Request) (claims interface{}, roleGranted *role.Role, shouldLogin bool, err error) {
			return nil, nil, true, nil
		}

		loginHandler, err := accesscontrol.NewLogin(lah)
		if err != nil {
			t.Fatalf("Error while creating new login handler, error: %v", err)
		}

		ctx := context.Background()
		r := httptest.NewRequest("GET", "http://foo.com", nil)
		rw := httptest.NewRecorder()

		err = loginHandler(ctx, rw, r)
		if err != nil {
			if !errors.Is(err, ErrNilRoleGranted) {
				t.Fatalf("Expected ErrNilRoleGranted while executing loginHandler with invalid lah")
			}
		} else {
			t.Fatalf("Expected error while executing loginHandler with invalid lah")
		}

	})
}
