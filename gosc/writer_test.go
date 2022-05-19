package gosc

import (
	"bufio"
	"bytes"
	"testing"
)

type fakePackage struct{}

func (f *fakePackage) GetType() PackageType {
	return "fake"
}

func Test_writePackage(t *testing.T) {
	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)

	t.Run("message", func(t *testing.T) {
		err := writePackage(&Message{
			Address:   "/info",
			Arguments: []any{},
		}, w)
		if err != nil {
			t.Errorf("expected no error but got: %v", err)
		}
	})
	t.Run("bundle", func(t *testing.T) {
		err := writePackage(&Bundle{}, w)
		if err == nil {
			t.Error("expected error but none given")
		}
	})
	t.Run("unknown", func(t *testing.T) {
		err := writePackage(&fakePackage{}, w)
		if err == nil {
			t.Error("expected error but none given")
		}
	})
}

func Test_writeArguments(t *testing.T) {
	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)
	err := writeArguments(w, []any{})
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	_ = w.Flush()
	if buf.Len() != 4 {
		t.Errorf("arguments was %d bytes but expected 4", buf.Len())
	}
}

func Test_writePaddedString(t *testing.T) {
	buf := bytes.Buffer{}
	w := bufio.NewWriter(&buf)
	data := "/info"
	err := writePaddedString(w, data)
	if err != nil {
		t.Errorf("expected no error but got: %v", err)
	}
	_ = w.Flush()
	if buf.Len() != 8 {
		t.Errorf("padded string was %d bytes but expected 8", buf.Len())
	}
}
