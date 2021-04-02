package accesscontrol

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corioders/gokit/errors"
	"github.com/corioders/gokit/web/middleware/accesscontrol/role"
)

func TestNewVerify(t *testing.T) {
	accesscontrol, err := New("TestNewVerify, accesscontrol, number 1", validAccesscontrolKey)
	if err != nil {
		t.Fatalf("Error while creating new accesscontrol, error: %v", err)
	}

	t.Run("ValidPermissionsNeeded", func(t *testing.T) {
		roleManager, err := role.NewManager("TestNewVerify, roleManager, number 1")
		if err != nil {
			t.Fatalf("Error while creating new roleManager, error: %v", err)
		}

		normalPermission, err := roleManager.NewPermission("TestNewVerify, permission, number 1")
		if err != nil {
			t.Fatalf("Error while creating new permission, error: %v", err)
		}

		_, err = accesscontrol.NewVerify([]*role.Permission{normalPermission})
		if err != nil {
			t.Fatalf("Expected no error while creating new verify middelware, but got error: %v", err)
		}
	})

	t.Run("NilPermissionsNeeded", func(t *testing.T) {
		_, err = accesscontrol.NewVerify(nil)
		if err != nil {
			if !errors.Is(err, ErrNilPermissionsNeeded) {
				t.Fatalf("Expected ErrNilPermissionsNeeded while creating new verify middelware with nil permissionsNeeded")
			}
		} else {
			t.Fatalf("Expected error while creating new verify middelware with nil permissionsNeeded")
		}
	})

	t.Run("ZeroPermissionsNeeded", func(t *testing.T) {
		_, err = accesscontrol.NewVerify([]*role.Permission{})
		if err != nil {
			if !errors.Is(err, ErrZeroPermissionsNeeded) {
				t.Fatalf("Expected ErrZeroPermissionsNeeded while creating new verify middelware with zero permissionsNeeded")
			}
		} else {
			t.Fatalf("Expected error while creating new verify middelware with zero permissionsNeeded")
		}
	})
}

func TestVerifyMiddelware(t *testing.T) {
	accesscontrol, err := New("TestVerifyMiddelware, accesscontrol, number 1", validAccesscontrolKey)
	if err != nil {
		t.Fatalf("Error while creating new accesscontrol, error: %v", err)
	}

	roleManager, err := role.NewManager("TestVerifyMiddelware, roleManager, number 1")
	if err != nil {
		t.Fatalf("Error while creating new roleManager, error: %v", err)
	}

	normalPermission, err := roleManager.NewPermission("TestVerifyMiddelware, permission, normal permission")
	if err != nil {
		t.Fatalf("Error while creating new permission, error: %v", err)
	}

	notNormalPermission, err := roleManager.NewPermission("TestVerifyMiddelware, permission, not normal permission")
	if err != nil {
		t.Fatalf("Error while creating new permission, error: %v", err)
	}

	normalRole, err := roleManager.NewRole("TestVerifyMiddelware, role, normal role", normalPermission)
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	validToken, err := newLoginToken(accesscontrol.tokenSigner, accesscontrol.tokenEncrypter, &internalClaims{
		UserClaims: nil,
		Role:       normalRole,
	})
	if err != nil {
		t.Fatalf("Error while creating new token, error: %v", err)
	}

	t.Run("Allowed", func(t *testing.T) {
		verifyMiddelware, err := accesscontrol.NewVerify([]*role.Permission{normalPermission})
		if err != nil {
			t.Fatalf("Error while creating new verify middelware, error: %v", err)
		}

		handlerExecuted := false
		verifyHandler := verifyMiddelware(func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
			handlerExecuted = true
			return nil
		})

		ctx := context.Background()
		r := httptest.NewRequest("GET", "http://foo.com", nil)
		r.AddCookie(&http.Cookie{
			Name:  accesscontrol.tokenCookieName,
			Value: validToken,
		})
		rw := httptest.NewRecorder()

		err = verifyHandler(ctx, rw, r)
		if err != nil {
			t.Fatalf("Error while executing verify middelware, error: %v", err)
		}

		if !handlerExecuted {
			t.Fatal("Expected handler to be executed when token is valid and has right permissions")
		}
	})

	t.Run("Disallowed", func(t *testing.T) {
		verifyMiddelware, err := accesscontrol.NewVerify([]*role.Permission{normalPermission, notNormalPermission})
		if err != nil {
			t.Fatalf("Error while creating new verify middelware, error: %v", err)
		}

		handlerExecuted := false
		verifyHandler := verifyMiddelware(func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
			handlerExecuted = true
			return nil
		})

		ctx := context.Background()
		r := httptest.NewRequest("GET", "http://foo.com", nil)
		r.AddCookie(&http.Cookie{
			Name:  accesscontrol.tokenCookieName,
			Value: validToken,
		})
		rw := httptest.NewRecorder()

		err = verifyHandler(ctx, rw, r)
		if err != nil {
			t.Fatalf("Error while executing verify middelware, error: %v", err)
		}

		if handlerExecuted {
			t.Fatal("Expected handler not to be executed when token is valid and has less permissions")
		}
	})

	t.Run("NoCookie", func(t *testing.T) {
		verifyMiddelware, err := accesscontrol.NewVerify([]*role.Permission{normalPermission})
		if err != nil {
			t.Fatalf("Error while creating new verify middelware, error: %v", err)
		}

		handlerExecuted := false
		verifyHandler := verifyMiddelware(func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
			handlerExecuted = true
			return nil
		})

		ctx := context.Background()
		r := httptest.NewRequest("GET", "http://foo.com", nil)
		rw := httptest.NewRecorder()

		err = verifyHandler(ctx, rw, r)
		if err != nil {
			t.Fatalf("Error while executing verify middelware, error: %v", err)
		}

		if handlerExecuted {
			t.Fatal("Expected handler not to be executed when token there is no token cookie")
		}

		if rw.Result().StatusCode != statusAccessDenied {
			t.Fatal("Expected status code to be statusAccessDenied when there is no token cookie")
		}
	})

	t.Run("NilRole", func(t *testing.T) {
		tokenWithNilRole, err := newLoginToken(accesscontrol.tokenSigner, accesscontrol.tokenEncrypter, &internalClaims{
			UserClaims: nil,
			Role:       nil,
		})
		if err != nil {
			t.Fatalf("Error while creating new token, error: %v", err)
		}

		verifyMiddelware, err := accesscontrol.NewVerify([]*role.Permission{normalPermission})
		if err != nil {
			t.Fatalf("Error while creating new verify middelware, error: %v", err)
		}

		verifyHandler := verifyMiddelware(func(_ context.Context, _ http.ResponseWriter, _ *http.Request) error {
			return nil
		})

		ctx := context.Background()
		r := httptest.NewRequest("GET", "http://foo.com", nil)
		r.AddCookie(&http.Cookie{
			Name:  accesscontrol.tokenCookieName,
			Value: tokenWithNilRole,
		})
		rw := httptest.NewRecorder()

		err = verifyHandler(ctx, rw, r)
		if err != nil {
			if !errors.Is(err, ErrGotNilRole) {
				t.Fatalf("Expected ErrGotNilRole while executing verifyHandler with token.claims.role=nil")
			}
		} else {
			t.Fatalf("Expected error while executing verifyHandler with token.claims.role=nil")
		}
	})

	t.Run("UserClaims", func(t *testing.T) {
		type claims struct {
			Foo string
			Bar string
		}

		tokenWithNilRole, err := newLoginToken(accesscontrol.tokenSigner, accesscontrol.tokenEncrypter, &internalClaims{
			UserClaims: &userClaims{
				value: claims{
					Foo: "foo",
					Bar: "bar",
				},
			},
			Role: normalRole,
		})
		if err != nil {
			t.Fatalf("Error while creating new token, error: %v", err)
		}

		verifyMiddelware, err := accesscontrol.NewVerify([]*role.Permission{normalPermission})
		if err != nil {
			t.Fatalf("Error while creating new verify middelware, error: %v", err)
		}

		verifyHandler := verifyMiddelware(func(ctx context.Context, _ http.ResponseWriter, _ *http.Request) error {
			c := claims{}
			err := ctx.Value(CtxKeyGetClaims).(GetClaims)(&c)
			if err != nil {
				t.Fatalf("Expected no error while calling GetClaims, but got error: %v", err)
			}

			if c.Foo != "foo" || c.Bar != "bar" {
				t.Fatal("Expected claims to be the same as provided one")
			}
			return nil
		})

		ctx := context.Background()
		r := httptest.NewRequest("GET", "http://foo.com", nil)
		r.AddCookie(&http.Cookie{
			Name:  accesscontrol.tokenCookieName,
			Value: tokenWithNilRole,
		})
		rw := httptest.NewRecorder()

		err = verifyHandler(ctx, rw, r)
		if err != nil {
			t.Fatalf("Error while executing verify middelware, error: %v", err)
		}
	})
}
