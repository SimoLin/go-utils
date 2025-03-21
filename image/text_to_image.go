package image

import (
	"bufio"
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/SimoLin/go-utils/hash"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

type TextDrawer struct {
	DPI      float64 // screen resolution in Dots Per Inch
	FontFile string  // filename of the ttf font
	Hinting  string  // none | full
	Size     float64 // font size in points
	Height   int     // image height in points
	Width    int     // image width in points
	Padding  int     // text left and right padding
	Chars    int     // chars displayed per line
	Spacing  float64 // line spacing
	Wonb     bool    // white text on a black background
}

func NewTextDrawer(width int, height int, chars int) (t *TextDrawer) {
	t = &TextDrawer{
		DPI:      240,
		FontFile: "InconsolataYahei.ttf", // 默认使用 InconsolataYahei 字体，请按需修改
		Hinting:  "none",
		Size:     12,
		Height:   height,
		Width:    width,
		Padding:  10,
		Chars:    chars,
		Spacing:  1.0,
		Wonb:     false,
	}
	return
}

var DefaultTextDrawer = NewTextDrawer(0, 0, 0)

func GetTextLength(text string) (length int) {
	length_string := len([]rune(text))
	length_bytes := len([]byte(text))
	return int((length_string + length_bytes) / 2)
}

func GetMaxTextLength(text []string) (max_text_length int, count_line int) {
	for _, line := range text {
		length_current_line := GetTextLength(line)
		if length_current_line > max_text_length {
			max_text_length = length_current_line
		}
	}
	count_line = len(text)
	return
}

func (t *TextDrawer) TextToImage(text []string) (rgba *image.RGBA, err error) {

	// 当宽高未指定时，自适应宽高
	if t.Width == 0 && t.Height == 0 {
		max_text_length, count_line := GetMaxTextLength(text)
		t.Width = int(t.Size*t.DPI/72*float64(max_text_length)/2) + t.Padding*2
		t.Height = int(t.Size*t.DPI/72*float64(count_line)) + t.Padding*2
	}

	// 当指定Chars数量时，自动修改单个字符大小
	if t.Chars > 0 {
		t.Size = float64(t.Width-t.Padding*2) / float64(t.Chars) * 72 / t.DPI
	}

	f, err := t.ReadFontFile()
	if err != nil {
		return
	}

	fg, bg := image.Black, image.White
	if t.Wonb {
		fg, bg = bg, fg
	}

	rgba = image.NewRGBA(image.Rect(0, 0, t.Width, t.Height))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{}, draw.Src)

	// Freetype context
	c := freetype.NewContext()
	c.SetDPI(t.DPI)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetFont(f)
	c.SetFontSize(t.Size)
	c.SetSrc(fg)
	switch t.Hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	opts := truetype.Options{}
	opts.Size = t.Size
	opts.DPI = t.DPI
	face := truetype.NewFace(f, &opts)

	// Calculate the widths and print to image
	pt := freetype.Pt(t.Padding, c.PointToFixed(t.Size).Round())
	newline := func() {
		pt.X = fixed.Int26_6(t.Padding) << 6
		pt.Y += c.PointToFixed(t.Size * t.Spacing)
	}

	for _, line := range text {
		for _, x := range line {
			w, _ := face.GlyphAdvance(x)
			if x == '\t' {
				x = ' '
			} else if f.Index(x) == 0 {
				continue
			} else if pt.X.Round()+w.Round() > t.Width-t.Padding {
				newline()
			}

			pt, err = c.DrawString(string(x), pt)
			if err != nil {
				log.Fatal(err)
			}
		}
		newline()
	}

	// t.SaveImageToFile(rgba, "out.png")

	return
}

func (t *TextDrawer) ReadFontFile() (f *truetype.Font, err error) {
	b, err := os.ReadFile(t.FontFile)
	if err != nil {
		log.Panic(err)
	}
	f, err = truetype.Parse(b)
	if err != nil {
		log.Panic(err)
	}
	return
}

func (t *TextDrawer) SaveImageToFile(rgba image.Image, fileName string) {
	// Save that RGBA image to disk.
	outFile, err := os.Create(fileName)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	err = b.Flush()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	fmt.Println("Wrote out.png OK.")
}

func (t *TextDrawer) GetImageBase64AndMD5(rgba image.Image) (image_base64 string, image_MD5 string) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, rgba, nil)
	if err != nil {
		return
	}
	image_bytes := buf.Bytes()
	image_base64 = hash.Base64Encode(string(image_bytes))
	image_MD5 = hash.MD5EncodeByte(image_bytes)
	return
}
