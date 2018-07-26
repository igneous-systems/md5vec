/*
Package md5vec implements an N-way parallel (up to 8-way) MD5 Write()
on Intel machines with AVX2.
*/
package md5vec

import (
	"crypto/md5"
	"fmt"
	"hash"
	"unsafe"
)

// WriteN writes the contents of N byte buffers to the
// hash state of N md5 hash digests. The lengths of all
// N buffers must be identical.
func WriteN(h []hash.Hash, p [][]byte) (n int, err error) {
	if len(h) != len(p) {
		return 0, fmt.Errorf("mismatched buffer count (%d) and digest count (%d)", len(h), len(p))
	}
	if len(h) == 0 {
		return
	}
	n = len(p[0])
	for i := range p {
		if n != len(p[i]) {
			return 0, fmt.Errorf("mismatched buf lengths: len(buf[%d])=%d, len(buf[0])=%d", i, len(p[i]), n)
		}
	}
	if n == 0 {
		return
	}
	if !hasAVX2 || len(h) == 1 {
		for i := range h {
			h[i].Write(p[i])
		}
		return
	}

	// Round n down to a multiple of the largest block size,
	// and store the true byte count in m. Trailing bytes are
	// written with md5.Write.
	m := n
	n &^= md5.BlockSize - 1

	for j := 0; j < len(h); {
		// choose at most 8 source buffers
		sh := h[j:]
		if len(sh) > 8 {
			sh = sh[:8]
		}
		sp := p[j:]
		if len(sp) > 8 {
			sp = sp[:8]
		}

		// each buffer must be pointed to with a 32 bit
		// displacement; one of the vector registers contains
		// the displacements from a base register to conduct
		// 8x32 bit loads.
		var bufs [8]int32

		// apparently the offsets must be positive and nonzero!
		// so pick a base register 4 bytes below the smallest
		// buf pointer.
		base := uintptr(unsafe.Pointer(&sp[0][0])) - 4
		for i := range sh {
			bp := uintptr(unsafe.Pointer(&sp[i][0])) - 4
			if bp < base {
				base = bp
			}
		}

		// condense up to 8 hash states at a time into
		// 4x8 uint32 vectors
		var s digest8
		for i := range sh {
			d := digestOf(&sh[i])
			if d.nx != 0 {
				return 0, fmt.Errorf("partial state of sub-Block unsupported")
			}
			s.v0[i] = d.s[0]
			s.v1[i] = d.s[1]
			s.v2[i] = d.s[2]
			s.v3[i] = d.s[3]

			diff := uintptr(unsafe.Pointer(&sp[i][0])) - base
			if int64(int32(diff)) != int64(diff) {
				return 0, fmt.Errorf("buffers exceed span of 32 bit displacement")
			}
			bufs[i] = int32(diff)
		}

		var cache cache8 // stack storage for block8 tmp state
		block8(&s.v0[0], base, &bufs[0], &cache[0], n)

		// burst hash state into constituent digests
		for i := range sh {
			d := digestOf(&sh[i])
			d.s[0] = s.v0[i]
			d.s[1] = s.v1[i]
			d.s[2] = s.v2[i]
			d.s[3] = s.v3[i]
			d.len += uint64(n)
		}
		j += len(sh)
	}
	// write trailing bytes
	if n != m {
		for i := range h {
			h[i].Write(p[i][n:])
		}
		n = m
	}
	return
}
