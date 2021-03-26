package rand

import "strings"

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index.
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits.
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits.
)

// String generates random string with length n, returned string is url-safe.
func (r *Rand) String(n int) (string, error) {
	sb := strings.Builder{}
	sb.Grow(n)

	cache, err := r.Uint64()
	if err != nil {
		return "", err
	}
	remain := letterIdxMax
	for i := n - 1; i >= 0; {
		if remain == 0 {
			cache, err = r.Uint64()
			if err != nil {
				return "", err
			}

			remain = letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String(), nil
}
