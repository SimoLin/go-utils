package test

import (
	"fmt"
	"testing"

	"github.com/SimoLin/go-utils/image"
)

func TestTextToImage(t *testing.T) {
	text := []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"测试bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb测试",
	}

	image.DefaultTextDrawer.TextToImage(text)
}

func TestGetTextLength(t *testing.T) {
	text := map[int]string{
		40: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		44: "测试bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	}
	for length_correct, line := range text {
		length := image.GetTextLength(line)
		fmt.Println(length)
		if length_correct != length {
			t.Error()
		}
	}
}
