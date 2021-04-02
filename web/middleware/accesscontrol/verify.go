package accesscontrol

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/corioders/gokit/errors"
	"github.com/corioders/gokit/web"
	"github.com/corioders/gokit/web/middleware/accesscontrol/role"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type verifyOptions struct{}
type verifyOption func(o *verifyOptions)

var (
	ErrNilPermissionsNeeded  = errors.New("PermissionsNeeded cannot be nil")
	ErrZeroPermissionsNeeded = errors.New("PermissionsNeeded cannot be of length zero because then everyone can access")

	ErrGotNilRole       = errors.New("Got nil role while reading token")
	ErrGettingNilClaims = errors.New("Calling GetClaims when provided claims are nil")
)

// NewVerify creates new verification midellware, user needs to have all permissions from permissionsNeeded to be allowed.
func (ac *Accesscontrol) NewVerify(permissionsNeeded []*role.Permission, options ...verifyOption) (web.Middleware, error) {
	if permissionsNeeded == nil {
		return nil, errors.WithStack(ErrNilPermissionsNeeded)
	}
	if len(permissionsNeeded) == 0 {
		return nil, errors.WithStack(ErrZeroPermissionsNeeded)
	}

	// Record stack, so finding call to NewVerify is easy.
	errGotNilRole := errors.WithStack(ErrGotNilRole)

	return func(handler web.Handler) web.Handler {
		return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			tokenCookie, err := r.Cookie(ac.tokenCookieName)
			if err != nil {
				if errors.Is(err, http.ErrNoCookie) {
					rw.WriteHeader(statusAccessDenied)
					return nil
				}
				return errors.WithStack(err)
			}

			tokenEncrypted, err := jwt.ParseSignedAndEncrypted(tokenCookie.Value)
			if err != nil {
				return errors.WithStack(err)
			}

			token, err := tokenEncrypted.Decrypt(ac.encrypterKey)
			if err != nil {
				if errors.Is(err, jose.ErrCryptoFailure) {
					rw.WriteHeader(statusAccessDenied)
					return nil
				}

				return errors.WithStack(err)
			}

			claims := internalClaims{}
			err = token.Claims(ac.singerKey, &claims)
			if err != nil {
				if errors.Is(err, jose.ErrCryptoFailure) {
					rw.WriteHeader(statusAccessDenied)
					return nil
				}

				return errors.WithStack(err)
			}

			if claims.Role == nil {
				rw.WriteHeader(statusAccessDenied)
				return errors.WithStack(errGotNilRole)
			}

			for _, p := range permissionsNeeded {
				if !claims.Role.IsAllowedTo(p) {
					rw.WriteHeader(statusAccessDenied)
					return nil
				}
			}

			ctx = context.WithValue(ctx, CtxKeyGetClaims, GetClaims(func(v interface{}) error {
				if claims.UserClaims != nil {
					return json.Unmarshal(claims.UserClaims.unmarshalData, v)
				}
				return errors.WithStack(ErrGettingNilClaims)
			}))
			return handler(ctx, rw, r)
		}
	}, nil
}
