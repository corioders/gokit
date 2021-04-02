package accesscontrol

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/corioders/gokit/crypto/hash"
	"github.com/corioders/gokit/errors"
	"github.com/corioders/gokit/math/rand"
	"github.com/corioders/gokit/web/middleware/accesscontrol/role"
	"gopkg.in/square/go-jose.v2"
)

type ctxKey int

// GetClaims runs json.Unmarshal(claimsData, v), so v must be a pointer
type GetClaims func(v interface{}) error

const (
	CtxKeyGetClaims ctxKey = iota
)

var statusAccessDenied = http.StatusForbidden

var accesscontrolNames = &sync.Map{}

type Key = []byte
type Accesscontrol struct {
	tokenCookieName string
	key             Key

	tokenSigner jose.Signer
	singerKey   Key

	tokenEncrypter jose.Encrypter
	encrypterKey   Key
}

type internalClaims struct {
	UserClaims *userClaims `json:"uc,omitempty"`
	Role       *role.Role  `json:"ro,omitempty"`
}

type userClaims struct {
	value         interface{}
	unmarshalData []byte
}

// UnmarshalJSON implements Unmarshaler interface.
func (uc *userClaims) UnmarshalJSON(data []byte) error {
	uc.unmarshalData = data
	return nil
}

// MarshalJSON implements Marshaler interface.
func (uc *userClaims) MarshalJSON() ([]byte, error) {
	return json.Marshal(uc.value)
}

var (
	ErrInvalidKeyLength = errors.New("Invalid key size, expected 64 bytes")
	ErrNameNonUnique    = errors.New("Name must be unique")
)

// New creates new Accesscontrol instance, name must be unique to every New call.
// Key must be 64 bytes in size, if not ErrInvalidKeyLength is returned.
// Name must be unique, if not ErrNameNonUnique is returned.
func New(name string, key Key) (*Accesscontrol, error) {
	if len(key) != 64 {
		return nil, errors.WithStack(ErrInvalidKeyLength)
	}

	_, ok := accesscontrolNames.Load(name)
	if ok {
		return nil, errors.WithMessage(ErrNameNonUnique, fmt.Sprintf(`name "%v" is not unique`, name))
	}

	singerKey := key[:32]
	tokenSinger, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.HS256, Key: singerKey}, &jose.SignerOptions{
		NonceSource: rand.NewCrypto(),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	encrypterKey := key[32:]
	tokenEncrypter, err := jose.NewEncrypter(jose.A256GCM, jose.Recipient{Algorithm: jose.A128GCMKW, Key: encrypterKey}, &jose.EncrypterOptions{
		Compression: jose.DEFLATE,
		ExtraHeaders: map[jose.HeaderKey]interface{}{
			jose.HeaderContentType: jose.ContentType("JWT"),
		},
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	accesscontrolNames.Store(name, struct{}{})
	return &Accesscontrol{
		key:             key,
		tokenCookieName: hash.Sha256Base64UrlSafe([]byte(name)),

		tokenSigner: tokenSinger,
		singerKey:   singerKey,

		tokenEncrypter: tokenEncrypter,
		encrypterKey:   encrypterKey,
	}, nil
}
