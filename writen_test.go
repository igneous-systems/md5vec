package md5vec

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"hash"
	"math/rand"
	"strings"
	"testing"
	"unsafe"
)

const (
	nbufs   = 8
	bufsize = 1024
)

func TestWriteNull(t *testing.T) {
	n, err := WriteN(nil, nil)
	if n != 0 || err != nil {
		t.Fatal(err)
	}
}

func TestWriteZero(t *testing.T) {
	n, err := WriteN([]hash.Hash{md5.New()}, [][]byte{nil})
	if n != 0 || err != nil {
		t.Fatal(err)
	}
}

func TestMismatchCount(t *testing.T) {
	n, err := WriteN(nil, [][]byte{nil})
	if err == nil || !strings.Contains(err.Error(), "mismatched buffer count") {
		t.Fatalf("expected mismatch, got: %v %v", n, err)
	}
}

func TestMismatchLength(t *testing.T) {
	n, err := WriteN([]hash.Hash{md5.New(), md5.New()}, [][]byte{[]byte("a"), []byte("aa")})
	if err == nil || !strings.Contains(err.Error(), "mismatched buf lengths") {
		t.Fatalf("expected mismatch, got: %v %v", n, err)
	}
}

func TestNonBlocksize(t *testing.T) {
	src := bytes.Repeat([]byte("a"), 99)
	h8 := md5.New()
	h := md5.New()
	WriteN([]hash.Hash{h8}, [][]byte{src})
	h.Write(src)
	if !bytes.Equal(h8.Sum(nil), h.Sum(nil)) {
		t.Fatal("compare")
	}
}

func TestPanicPartialWrite(t *testing.T) {
	// N/A for scalar md5 code
	if !hasAVX2 {
		t.Skip("avx unsupported")
	}

	h1 := md5.New()
	h2 := md5.New()
	h1.Write([]byte("hi"))
	h2.Write([]byte("hi"))
	sixty4 := bytes.Repeat([]byte("a"), 64)
	n, err := WriteN([]hash.Hash{h1, h2}, [][]byte{sixty4, sixty4})
	if err == nil || !strings.Contains(err.Error(), "partial state") {
		t.Fatalf("expected partial state, got: %v %v", n, err)
	}
}

func TestWriteN(t *testing.T) {
	for i := 1; i < nbufs; i++ {
		if err := testWriteN(t, i); err != nil {
			t.Errorf("nbufs %d: %v", i, err)
		}
	}
}

func testWriteN(t *testing.T, nb int) error {
	p := make([][]byte, nb)
	for i := range p {
		p[i] = make([]byte, bufsize)
		for j := range p[i] {
			p[i][j] = byte(rand.Int())
		}
	}

	mh := make([]hash.Hash, nb)
	for i := range mh {
		mh[i] = md5.New()
	}

	mh8 := make([]hash.Hash, nb)
	for i := range mh {
		mh8[i] = md5.New()
	}

	for i := range mh {
		mh[i].Write(p[i])
	}

	n, err := WriteN(mh8, p)
	if err != nil {
		return err
	}
	if n != bufsize {
		return fmt.Errorf("n(%d) != bufsize(%d)", n, bufsize)
	}

	var badslots []int
	var offsets []int
	for i := range p {
		exp := mh[i].Sum(nil)
		got := mh8[i].Sum(nil)
		if !bytes.Equal(exp, got) {
			badslots = append(badslots, i)
			offsets = append(offsets, int(uintptr(unsafe.Pointer(&p[i][0]))-uintptr(unsafe.Pointer(&p[0][0]))))
		}
	}
	if badslots != nil {
		for i := range p {
			t.Logf("%p %d", &p[i][0], int(uintptr(unsafe.Pointer(&p[i][0]))-uintptr(unsafe.Pointer(&p[0][0]))))
		}
		return fmt.Errorf("bad compares at %v %v", badslots, offsets)
	}
	return nil
}

func TestPartialN(t *testing.T) {
	p := make([][]byte, 2)
	mh := make([]hash.Hash, 2)
	for i := range p {
		p[i] = make([]byte, bufsize)
		mh[i] = md5.New()
		for j := range p[i] {
			p[i][j] = byte(rand.Int())
		}
	}

	for i := 0; i < bufsize; i += md5.BlockSize {
		q := make([][]byte, 2)
		for j := range q {
			q[j] = p[j][i : i+md5.BlockSize]
		}
		WriteN(mh, q)
	}

	for i := range mh {
		h := md5.New()
		h.Write(p[i])
		if !bytes.Equal(mh[i].Sum(nil), h.Sum(nil)) {
			t.Fatal("compare")
		}
	}
}

func BenchmarkMd5(b *testing.B) {
	p := make([][]byte, 8)
	for i := range p {
		p[i] = make([]byte, 1024*bufsize)
		for j := range p[i] {
			p[i][j] = byte(rand.Int())
		}
	}
	b.ResetTimer()
	b.SetBytes(1024 * 8 * bufsize)
	for i := 0; i < b.N; i++ {
		for j := range p {
			h := md5.New()
			h.Write(p[j])
			h.Sum(nil)
		}
	}
}

func BenchmarkMd5by8(b *testing.B) {
	xbuf := make([]byte, 8*1024*bufsize)
	p := make([][]byte, 8)
	for i := range p {
		p[i] = xbuf[i*1024*bufsize : (i*1024*bufsize)+1024*bufsize]
		for j := range p[i] {
			p[i][j] = byte(rand.Int())
		}
	}
	b.ResetTimer()
	b.SetBytes(8 * 1024 * bufsize)
	for i := 0; i < b.N; i++ {
		h := make([]hash.Hash, len(p))
		for j := range h {
			h[j] = md5.New()
		}
		WriteN(h, p)
		for j := range h {
			h[j].Sum(nil)
		}
	}
}
