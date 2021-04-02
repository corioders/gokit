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

type loginOptions struct{}
type loginOption func(o *loginOptions)
type LoginAcceptHandler func(ctx context.Context, r *http.Request) (claims interface{}, roleGranted *role.Role, shouldLogin bool, err error)

var (
	ErrNilLoginAcceptHandler = errors.New("Login accept handler (lah) cannot be nil")
	ErrNilRoleGranted        = errors.New("Login accept handler (lah) returned roleGranted=nil and shouldLogin=true, this is incorrect")
)

func (ac *Accesscontrol) NewLogin(lah LoginAcceptHandler, options ...loginOption) (web.Handler, error) {
	if lah == nil {
		return nil, errors.WithStack(ErrNilLoginAcceptHandler)
	}

	// Add stack when NewLogin is called so search for invalid lah is easy.
	errNilRoleGranted := errors.WithStack(ErrNilRoleGranted)

	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		claims, roleGranted, shouldLogin, err := lah(ctx, r)
		if err != nil {
			return errors.WithStack(err)
		}
		if !shouldLogin {
			rw.WriteHeader(statusAccessDenied)
			return nil
		}
		if roleGranted == nil {
			// errNilRoleGranted is tied to lah so we want to make search for invalid lah as easy as possible.
			return errNilRoleGranted
		}

		token, err := newLoginToken(ac.tokenSigner, ac.tokenEncrypter, &internalClaims{
			UserClaims: &userClaims{value: claims},
			Role:       roleGranted,
		})
		if err != nil {
			return errors.WithStack(err)
		}

		http.SetCookie(rw, &http.Cookie{
			Name:  ac.tokenCookieName,
			Value: token,
		})

		return nil
	}, nil
}

func newLoginToken(signer jose.Signer, encrypter jose.Encrypter, claims *internalClaims) (string, error) {
	return jwt.SignedAndEncrypted(signer, encrypter).Claims(claims).CompactSerialize()
}
