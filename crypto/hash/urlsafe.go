package hash

import (
	"crypto/sha256"
	"encoding/base64"
)

func Sha256Base64UrlSafe(data []byte) string {
	hash := sha256.Sum256(data)
	return base64.RawURLEncoding.EncodeToString(hash[:])
}
