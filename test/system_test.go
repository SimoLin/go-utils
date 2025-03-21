package test

import (
	"fmt"
	"testing"

	"github.com/SimoLin/go-utils/system"
)

func TestGetBoardInfo(t *testing.T) {
	board_model, err := system.GetBoardInfo()
	if err != nil {
		t.Failed()
	}
	fmt.Println(board_model)
}
