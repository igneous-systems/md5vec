// +build amd64

package md5vec

//go:noescape
func setAVX2()

//go:noescape
func block8(state *uint32, base uintptr, bufs *int32, cache *byte, n int)

func init() { setAVX2() }
