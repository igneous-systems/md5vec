package md5vec

import (
	"hash"
	"unsafe"
)

// Spoof the MD5 digest here as it is unexported from crypto/md5.
// An ugly hack, but a reflection/verification step is included
// in the unit test.

type digest struct {
	s   [4]uint32
	x   [64]byte
	nx  int
	len uint64
}

// digestOf returns the data portion of a runtime.iface as a *digest; assumes
// struct equivalence between digest and md5.digest (see the unit test)
func digestOf(hp *hash.Hash) *digest {
	iface := (*struct{ tab, data unsafe.Pointer })(unsafe.Pointer(hp))
	return (*digest)(iface.data)
}
