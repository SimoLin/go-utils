package text_drawer

import (
	"bufio"
	"bytes"
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
	FontSize float64 // font size in points
	Height   int     // image height in points
	Width    int     // image width in points
	Padding  int     // text left and right padding
	Chars    int     // chars displayed per line
	Spacing  float64 // line spacing
	Wonb     bool    // white text on a black background
}

type OptionFunc func(*TextDrawer)

func initOptions(options ...OptionFunc) *TextDrawer {
	opts := &TextDrawer{
		DPI:      240,
		FontFile: "InconsolataYahei.ttf", // 默认使用 InconsolataYahei 字体，请按需修改
		Hinting:  "none",
		FontSize: 12,
		Height:   0,
		Width:    0,
		Padding:  10,
		Chars:    0,
		Spacing:  1.0,
		Wonb:     false,
	}
	for _, option := range options {
		option(opts)
	}
	return opts
}

func WithOptions(t TextDrawer) OptionFunc {
	return func(opts *TextDrawer) {
		*opts = t
	}
}

func WithDPI(f float64) OptionFunc {
	return func(opts *TextDrawer) {
		opts.DPI = f
	}
}

func WithFontFile(s string) OptionFunc {
	return func(opts *TextDrawer) {
		opts.FontFile = s
	}
}

func WithHingting(s string) OptionFunc {
	return func(opts *TextDrawer) {
		opts.Hinting = s
	}
}

func WithFontSize(f float64) OptionFunc {
	return func(opts *TextDrawer) {
		opts.FontSize = f
	}
}

func WithHeight(i int) OptionFunc {
	return func(opts *TextDrawer) {
		opts.Height = i
	}
}

func WithWidth(i int) OptionFunc {
	return func(opts *TextDrawer) {
		opts.Width = i
	}
}

func WithPadding(i int) OptionFunc {
	return func(opts *TextDrawer) {
		opts.Padding = i
	}
}

func WithChars(i int) OptionFunc {
	return func(opts *TextDrawer) {
		opts.Chars = i
	}
}

func WithSpacing(f float64) OptionFunc {
	return func(opts *TextDrawer) {
		opts.Spacing = f
	}
}

func WithWonb(b bool) OptionFunc {
	return func(opts *TextDrawer) {
		opts.Wonb = b
	}
}

func New(options ...OptionFunc) *TextDrawer {
	opts := initOptions(options...)
	return opts
}

func (t *TextDrawer) TextToImage(text []string) (rgba *image.RGBA, err error) {

	// 当宽高未指定时，自适应宽高
	if t.Width == 0 && t.Height == 0 {
		max_text_length, count_line := GetMaxTextLength(text)
		t.Width = int(t.FontSize*t.DPI/72*float64(max_text_length)/2) + t.Padding*2
		t.Height = int(t.FontSize*t.DPI/72*float64(count_line)) + t.Padding*2
	}

	// 当指定Chars数量时，自动修改单个字符大小
	if t.Chars > 0 {
		t.FontSize = float64(t.Width-t.Padding*2) / float64(t.Chars) * 72 / t.DPI
	}

	f, err := ReadFontFile(t.FontFile)
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
	c.SetFontSize(t.FontSize)
	c.SetSrc(fg)
	switch t.Hinting {
	default:
		c.SetHinting(font.HintingNone)
	case "full":
		c.SetHinting(font.HintingFull)
	}

	opts := truetype.Options{}
	opts.Size = t.FontSize
	opts.DPI = t.DPI
	face := truetype.NewFace(f, &opts)

	// Calculate the widths and print to image
	pt := freetype.Pt(t.Padding, c.PointToFixed(t.FontSize).Round())
	newline := func() {
		pt.X = fixed.Int26_6(t.Padding) << 6
		pt.Y += c.PointToFixed(t.FontSize * t.Spacing)
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

	return
}

func ReadFontFile(font_file_path string) (f *truetype.Font, err error) {
	b, err := os.ReadFile(font_file_path)
	if err != nil {
		return
	}
	f, err = truetype.Parse(b)
	if err != nil {
		return
	}
	return
}

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

func ImageToByte(rgba image.Image) (image_bytes []byte) {
	image_bytes = make([]byte, 0)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, rgba, nil)
	if err != nil {
		return
	}
	return buf.Bytes()
}

func SaveImageToFile(rgba image.Image, file_path string) {
	outFile, err := os.Create(file_path)
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
}

func GetImageBase64(rgba image.Image) (image_base64 string) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, rgba, nil)
	if err != nil {
		return
	}
	image_bytes := buf.Bytes()
	image_base64 = hash.Base64Encode(string(image_bytes))
	return
}

func GetImageMD5(rgba image.Image) (image_MD5 string) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, rgba, nil)
	if err != nil {
		return
	}
	image_bytes := buf.Bytes()
	image_MD5 = hash.MD5EncodeByte(image_bytes)
	return
}
