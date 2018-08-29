# md5vec
AVX-accelerated 8-way parallel md5

This package contains the Intel AVX-specific code for a parallel md5 kernel.

The function which this package exports has the signature:

    func WriteN(h []hash.Hash, p [][]byte) (n int, err error)
    
which provides an up-to-8-way parallel invocation of md5 rounds 
on the given byte slices. (by analogy to `hash.Write`)

Note: all slices must be of the same length, and `len(p)` must be the
same as `len(h)`. Furthermore, `h` must represent a collection of `crypto/md5`
hash digests, as created by `md5.New()`. Therefore, after calling
`WriteN()` any trailers may be written directly via `hash.Write()` and
`hash.Sum()` may be used to compute the final md5 checksum for each
digest.

An associated "server" will be published to provide transparent
parallelism via a drop-in replacement for the `hash.Hash` interface.
