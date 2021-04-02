package integration_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/corioders/gokit/math/rand"
	"github.com/corioders/gokit/web/middleware/accesscontrol"
	"github.com/corioders/gokit/web/middleware/accesscontrol/role"
)

func TestAccesscontrol(t *testing.T) {
	key := make(accesscontrol.Key, 64)
	_, err := io.ReadFull(rand.NewMath(0), key)
	if err != nil {
		t.Fatalf("Error while reading new key, error: %v", err)
	}

	ac, err := accesscontrol.New("TestAccesscontrol, accesscontrol, number 1", key)
	if err != nil {
		t.Fatalf("Error while creating new role manager, error: %v", err)
	}

	rm, err := role.NewManager("TestAccesscontrol, roleManager, number 1")
	if err != nil {
		t.Fatalf("Error while creating new role manager, error: %v", err)
	}

	normalPermission, err := rm.NewPermission("TestAccesscontrol, permission, normal permission")
	if err != nil {
		t.Fatalf("Error while creating new permission, error: %v", err)
	}

	normalRole, err := rm.NewRole("TestAccesscontrol, role, normal role", normalPermission)
	if err != nil {
		t.Fatalf("Error while creating new role, error: %v", err)
	}

	const userClaims = "test claims"
	lah := func(ctx context.Context, r *http.Request) (claims interface{}, roleGranted *role.Role, shouldLogin bool, err error) {
		return userClaims, normalRole, true, nil
	}

	loginHandler, err := ac.NewLogin(lah)
	if err != nil {
		t.Fatalf("Error while creating new loginHandler, error: %v", err)
	}

	verifyMiddelware, err := ac.NewVerify([]*role.Permission{normalPermission})
	if err != nil {
		t.Fatalf("Error while creating new verifyMiddelware, error: %v", err)
	}

	protectedHandler := verifyMiddelware(func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
		c := ""
		err := ctx.Value(accesscontrol.CtxKeyGetClaims).(accesscontrol.GetClaims)(&c)
		if err != nil {
			return err
		}

		rw.Write([]byte(c))
		return nil
	})

	loginCtx := context.Background()

	loginRequest := httptest.NewRequest("GET", "http://foo.com", nil)
	loginResponseWriter := httptest.NewRecorder()
	err = loginHandler(loginCtx, loginResponseWriter, loginRequest)
	if err != nil {
		t.Fatalf("Error while executing loginHandler, error: %v", err)	
	}

	loginResult := loginResponseWriter.Result()
	tokenCookie := loginResult.Cookies()[0]

	protectedRequest := httptest.NewRequest("GET", "http://foo.com", nil)
	protectedRequest.AddCookie(tokenCookie)
	protectedResponseWriter := httptest.NewRecorder()
	protectedCtx := context.Background()
	err = protectedHandler(protectedCtx, protectedResponseWriter, protectedRequest)
	if err != nil {
		t.Fatalf("Error while executing protectedHandler, error: %v", err)	
	}

	result := protectedResponseWriter.Body.String()

	if result != userClaims {
		t.Fatal("Expected result to be equal to claims because protectedHandler writes claims to responseWriter")
	}
}
