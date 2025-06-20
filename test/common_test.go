package test

import (
	"fmt"
	"testing"

	"github.com/SimoLin/go-utils/common"
)

func TestSliceToJsonUint(t *testing.T) {
	s := "[1,2,3]"
	r, _ := common.StringToSliceUint(s)
	fmt.Println(r)
}

func TestSliceToSliceUint(t *testing.T) {
	s := "[1,2,3]"
	s1, _ := common.StringToSliceAny(s)
	r := common.SliceAnyToSliceUint(s1)
	fmt.Println(r)
}
