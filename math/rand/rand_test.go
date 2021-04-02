package rand

import (
	"bytes"
	"io"
	mathrand "math/rand"
	"testing"
	"unsafe"

	"github.com/corioders/gokit/errors"
)

func TestNew(t *testing.T) {
	t.Run("reader", func(t *testing.T) {
		sourceBytes := [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
		sourceBytesReader := bytes.NewReader(sourceBytes[:])
		randBytes := New(sourceBytesReader)

		if randBytes.source != sourceBytesReader {
			t.Fatal("Expected randBytes.source to equal sourceBytesReader")
		}

		v, err := randBytes.Uint64()
		if err != nil {
			t.Fatal("Error while reading Uint64", err)
		}

		vBytes := *(*[8]byte)(unsafe.Pointer(&v))
		if len(vBytes) != len(sourceBytes) {
			t.Fatal("Expected length of vBytes to equal sourceBytes, this is probably test error, not implementation bug")
		}

		for i := 0; i < len(vBytes); i++ {
			if vBytes[i] != sourceBytes[i] {
				t.Fatalf("Expected vBytes[%v] to equal sourceBytes[%v]", i, i)
			}
		}

		_, err = randBytes.Uint64()
		if !errors.Is(err, io.EOF) {
			t.Fatal("Expected error after source depletes to be io.EOF, got", err)
		}
	})

	t.Run("math", func(t *testing.T) {
		randMath := NewMath(0)

		if _, ok := randMath.source.(*mathrand.Rand); !ok {
			t.Fatal("Expected randMath.source to be of type math/rand.*Rand")
		}

		v, err := randMath.Uint64()
		if err != nil {
			t.Fatal("Error while reading Uint64", err)
		}

		sourceMath := make([]byte, 8)
		testSourceMathReader := mathrand.New(mathrand.NewSource(0))
		_, err = io.ReadFull(testSourceMathReader, sourceMath)
		if err != nil {
			t.Fatal("Error while reading from testSourceMathReader")
		}

		vBytes := *(*[8]byte)(unsafe.Pointer(&v))
		if len(vBytes) != len(sourceMath) {
			t.Fatal("Expected length of vBytes to equal sourceMath, this is probably test error, not implementation bug")
		}

		for i := 0; i < len(vBytes); i++ {
			if vBytes[i] != sourceMath[i] {
				t.Fatalf("Expected vBytes[%v] to equal sourceMath[%v]", i, i)
			}
		}
	})

}
