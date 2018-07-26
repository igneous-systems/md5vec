package md5vec

import (
	"crypto/md5"
	"reflect"
	"testing"
)

// verify md5 digest and cooked digest are reflection-equivalent
func TestMd5Interface(t *testing.T) {
	tp := reflect.TypeOf(md5.New()).Elem()
	dp := reflect.TypeOf(digest{})
	if tp.NumField() != dp.NumField() {
		t.Fatalf("unequal field count %d != %d", tp.NumField(), dp.NumField())
	}
	for i := 0; i < dp.NumField(); i++ {
		if tp.Field(i).Name != dp.Field(i).Name {
			t.Fatal(i, "names differ")
		}
		if !tp.Field(i).Type.AssignableTo(dp.Field(i).Type) {
			t.Fatal(i, "field not assignable")
		}
	}
}
