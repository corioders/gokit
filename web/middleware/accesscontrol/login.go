package accesscontrol

import (
	"context"
	"net/http"

	"github.com/corioders/gokit/web"
	"github.com/corioders/gokit/web/middleware/accesscontrol/role"
	"gopkg.in/square/go-jose.v2"
	"gopkg.in/square/go-jose.v2/jwt"
)

type loginOptions struct{}
type loginOption func(o *loginOptions)
type LoginAcceptHandler func(ctx context.Context, r *http.Request) (claims interface{}, roleGranted *role.Role, shouldLogin bool, err error)

func (ac *Accesscontrol) NewLogin(lah LoginAcceptHandler, options ...loginOption) web.Handler {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		claims, roleGranted, shouldLogin, err := lah(ctx, r)
		if err != nil {
			rw.WriteHeader(accessDeniedStatusCode)
			return err
		}
		if !shouldLogin {
			rw.WriteHeader(accessDeniedStatusCode)
			return nil
		}

		token, err := newLoginToken(ac.tokenSigner, ac.tokenEncrypter, &internalClaims{
			UserClaims: claims,
			Role:       roleGranted,
		})
		if err != nil {
			return err
		}

		http.SetCookie(rw, &http.Cookie{
			Name:  ac.tokenCookieName,
			Value: token,
		})

		return nil
	}
}

func newLoginToken(signer jose.Signer, encrypter jose.Encrypter, claims *internalClaims) (string, error) {
	return jwt.SignedAndEncrypted(signer, encrypter).Claims(claims).CompactSerialize()
}
