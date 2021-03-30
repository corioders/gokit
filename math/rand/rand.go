package rand

import (
	cryptorrand "crypto/rand"
	"io"
	mathrand "math/rand"
	"unsafe"

	"github.com/corioders/gokit/errors"
)

type Rand struct {
	source io.Reader
}

func New(source io.Reader) *Rand {
	return &Rand{source}
}

// NewCrypto is short for New("crypto/rand".Reader)
func NewCrypto() *Rand {
	return &Rand{cryptorrand.Reader}
}

func NewFromMath(source mathrand.Source) *Rand {
	return &Rand{mathrand.New(source)}
}

func (r *Rand) Read(p []byte) (n int, err error) {
	return r.source.Read(p)
}

var (
	ErrInvalidArgument = errors.New("invalid argument")
)

// Uint64 returns random number.
func (r *Rand) Uint64() (uint64, error) {
	bytes := [8]byte{}
	_, err := io.ReadFull(r.source, bytes[:])
	if err != nil {
		return 0, err
	}

	return *(*uint64)(unsafe.Pointer(&bytes)), nil
}

// Uint64n returns random number in range [0,n).
func (r *Rand) Uint64n(n uint64) (uint64, error) {
	v, err := r.Uint64()
	return v % n, err
}

// Uint64m returns random number in range [min,max).
// It returns error if min >= max.
func (r *Rand) Uint64m(min, max uint64) (uint64, error) {
	if min >= max {
		return 0, errors.WithMessage(ErrInvalidArgument, "to Uint64m, minx >= max")
	}

	v, err := r.Uint64n(max - min)
	return v + min, err
}

// Uint32 returns random number.
func (r *Rand) Uint32() (uint32, error) {
	bytes := [4]byte{}
	_, err := io.ReadFull(r.source, bytes[:])
	if err != nil {
		return 0, err
	}

	return *(*uint32)(unsafe.Pointer(&bytes)), nil
}

// Uint32n returns random number in range [0,n).
func (r *Rand) Uint32n(n uint32) (uint32, error) {
	v, err := r.Uint32()
	return v % n, err
}

// Uint64m returns random number in range [min,max).
// It returns error if min >= max.
func (r *Rand) Uint32m(min, max uint32) (uint32, error) {
	if min >= max {
		return 0, errors.WithMessage(ErrInvalidArgument, "to Uint32m, minx >= max")
	}

	v, err := r.Uint32n(max - min)
	return v + min, err
}
