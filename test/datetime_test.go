package test

import (
	"fmt"
	"testing"

	"github.com/SimoLin/go-utils/datetime"
)

func TestGetNowStringSimple(t *testing.T) {
	result := datetime.GetNowStringSimple()
	fmt.Println(result)
}
