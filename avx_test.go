package md5vec

import "testing"

func TestAvx(t *testing.T) {
	t.Log("you have avx2:", hasAVX2)
}
