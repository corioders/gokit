package rand

import (
	"strconv"
	"time"
)

func (r *Rand) Nonce() (string, error) {
	v, err := r.String(64)
	if err != nil {
		return "", err
	}

	return v + strconv.FormatInt(time.Now().Unix(), 10), nil
}
