package accesscontrol

import (
	"context"
	"net/http"

	"github.com/corioders/gokit/errors"
	"github.com/corioders/gokit/web"
	"github.com/corioders/gokit/web/middleware/accesscontrol/role"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type verifyOptions struct{}
type verifyOption func(o *verifyOptions)

func (ac *Accesscontrol) NewVerify(permissionsNeeded []*role.Permission, options ...verifyOption) web.Middleware {
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
					rw.WriteHeader(accessDeniedStatusCode)
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
					rw.WriteHeader(accessDeniedStatusCode)
					return nil
				}

				return errors.WithStack(err)
			}

			claims := internalClaims{}
			err = token.Claims(ac.singerKey, &claims)
			if err != nil {
				if errors.Is(err, jose.ErrCryptoFailure) {
					rw.WriteHeader(accessDeniedStatusCode)
					return nil
				}

				return errors.WithStack(err)
			}

			for _, p := range permissionsNeeded {
				if !claims.Role.IsAllowedTo(p) {
					rw.WriteHeader(accessDeniedStatusCode)
					return nil
				}
			}

			ctx = context.WithValue(ctx, CtxKeyClaims, claims.UserClaims)
			return handler(ctx, rw, r)
		}
	}
}
