// +build !amd64

package md5vec

// dummy declarations for AVX2 support for non-intel platforms

func block8(state *uint32, base uintptr, bufs *int32, cache *byte, n int) {
	panic("this shouldn't happen -- hasAVX2 returned true?")
}
