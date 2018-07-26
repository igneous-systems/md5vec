package md5vec

import (
	"crypto/md5"
	"fmt"
	"hash"
	"testing"
)

// 8 golden hashes from crypto/md5/md5_test.go longer than md5.BlockSize
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
var golden = [nbufs]struct{ out, in string }{
	{"637d2fe925c07c113800509964fb0e06", "For every action there is an equal and opposite government program."},
	{"de3a4d2fd6c73ec2db2abad23b444281", "There is no reason for any individual to have a computer in their home. -Ken Olsen, 1977"},
	{"acf203f997e2cf74ea3aff86985aefaf", "It's a tiny change to the code and not completely disgusting. - Bob Manchek"},
	{"cdf7ab6c1fd49bd9933c43f3ea5af185", "Give me a rock, paper and scissors and I will move the world.  CCFestoon"},
	{"277cbe255686b48dd7e8f389394d9299", "It's well we cannot hear the screams/That we create in others' dreams."},
	{"fd3fb0a7ffb8af16603f3d3af98f8e1f", "You remind me of a TV show, but that's all right: I watch it anyway."},
	{"63eb3a2f466410104731c4b037600110", "Even if I could be Shakespeare, I think I should still choose to be Faraday. - A. Huxley"},
	{"72c2ed7592debca1c90fc0100f931a2f", "The fugacity of a constituent in a mixture of gases at a given temperature is proportional to its mole fraction.  Lewis-Randall Rule"},
}

// Test the 8 golden hashes against an 8-way md5. Since each golden string is < 2*md5.BlockSize
// block8() will only run with one iteration; the tests in writen_test.go test loops of block8().
func TestGolden(t *testing.T) {
	h := make([]hash.Hash, nbufs)
	p := make([][]byte, nbufs)

	for i := range h {
		h[i] = md5.New()
	}
	for i := range p {
		p[i] = []byte(golden[i].in[:md5.BlockSize])
	}

	n, err := WriteN(h, p)
	if err != nil {
		t.Fatal(err)
	}
	if n != md5.BlockSize {
		t.Fatal(n)
	}

	for i := range h {
		h[i].Write([]byte(golden[i].in[md5.BlockSize:]))
	}
	for i := range h {
		if fmt.Sprintf("%x", h[i].Sum(nil)) != golden[i].out {
			t.Fatal("compare", i)
		}
	}
}
