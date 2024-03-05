package utils

import (
	"reflect"
	"unsafe"

	"github.com/njuwelkin/games/pal/mkf"
)

func WordArray(p unsafe.Pointer, size uintptr) []mkf.WORD {
	var w mkf.WORD
	l := size / unsafe.Sizeof(w)
	ret := []mkf.WORD{}
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	sh.Data = uintptr(p)
	sh.Len = int(l)
	sh.Cap = int(l)
	return ret
}
