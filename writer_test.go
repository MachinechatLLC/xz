package xz

import (
	"bytes"
	"io"
	"math/rand"
	"os"
	"testing"

	"github.com/ulikunitz/xz/randtxt"
)

func TestWriter(t *testing.T) {
	const text = "The quick brown fox jumps over the lazy dog."
	var buf bytes.Buffer
	w := NewWriter(&buf)
	n, err := io.WriteString(w, text)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}
	if n != len(text) {
		t.Fatalf("Writestring wrote %d bytes; want %d", n, len(text))
	}
	if err = w.Close(); err != nil {
		t.Fatalf("w.Close error %s", err)
	}
	var out bytes.Buffer
	r, err := NewReader(&buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	if _, err = io.Copy(&out, r); err != nil {
		t.Fatalf("io.Copy error %s", err)
	}
	s := out.String()
	if s != text {
		t.Fatalf("reader decompressed to %q; want %q", s, text)
	}
}

func TestWriter2(t *testing.T) {
	const txtlen = 59
	var buf bytes.Buffer
	io.CopyN(&buf, randtxt.NewReader(rand.NewSource(41)), txtlen)
	txt := buf.String()

	buf.Reset()
	w := NewWriter(&buf)
	n, err := io.WriteString(w, txt)
	if err != nil {
		t.Fatalf("WriteString error %s", err)
	}
	if n != len(txt) {
		t.Fatalf("WriteString wrote %d bytes; want %d", n, len(txt))
	}
	if err = w.Close(); err != nil {
		t.Fatalf("Close error %s", err)
	}
	f, _ := os.Create("test.xz")
	f.Write(buf.Bytes())
	f.Close()
	t.Logf("buf.Len() %d", buf.Len())
	r, err := NewReader(&buf)
	if err != nil {
		t.Fatalf("NewReader error %s", err)
	}
	var out bytes.Buffer
	k, err := io.Copy(&out, r)
	if err != nil {
		t.Fatalf("Decompressing copy error %s after %d bytes", err, n)
	}
	if k != txtlen {
		t.Fatalf("Decompression data length %d; want %d", k, txtlen)
	}
	if txt != out.String() {
		t.Fatal("decompressed data differes from original")
	}
}
