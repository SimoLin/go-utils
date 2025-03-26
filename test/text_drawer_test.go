package test

import (
	"fmt"
	"testing"

	"github.com/SimoLin/go-utils/text_drawer"
)

func TestTextToImage(t *testing.T) {
	text := []string{
		"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"测试bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb测试",
	}
	rgb, err := text_drawer.New(
		text_drawer.WithFontFile("InconsolataYahei.ttf"),
	).TextToImage(text)
	if err != nil {
		t.Fatal(err)
	}
	text_drawer.SaveImageToFile(rgb, "test.png")
	image_bytes := text_drawer.ImageToByte(rgb)
	image_base64 := text_drawer.GetImageBase64(rgb)
	image_md5 := text_drawer.GetImageMD5(rgb)
	fmt.Println(image_bytes)
	fmt.Println(image_base64)
	fmt.Println(image_md5)
}

func TestGetTextLength(t *testing.T) {
	text := map[int]string{
		40: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		44: "测试bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
	}
	for length_correct, line := range text {
		length := text_drawer.GetTextLength(line)
		fmt.Println(length)
		if length_correct != length {
			t.Error()
		}
	}
}
