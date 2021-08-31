package asciimage

import (
	"image"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/termenv"
	"github.com/nfnt/resize"
)

// Model is a struct that contains all the properties of the ascii image.
type Model struct {
	Image   image.Image
	Content string
	Height  int
	Width   int
}

// ImageToString converts an image to a string representation of an image.
func ImageToString(width, height uint, img image.Image) (string, error) {
	img = resize.Thumbnail(width, height*2-4, img, resize.Lanczos3)
	b := img.Bounds()
	w := b.Max.X
	h := b.Max.Y
	p := termenv.ColorProfile()
	str := strings.Builder{}
	for y := 0; y < h; y += 2 {
		for x := w; x < int(width); x = x + 2 {
			str.WriteString(" ")
		}
		for x := 0; x < w; x++ {
			c1, _ := colorful.MakeColor(img.At(x, y))
			color1 := p.Color(c1.Hex())
			c2, _ := colorful.MakeColor(img.At(x, y+1))
			color2 := p.Color(c2.Hex())
			str.WriteString(termenv.String("▀").
				Foreground(color1).
				Background(color2).
				String())
		}
		str.WriteString("\n")
	}
	return str.String(), nil
}

// SetContent sets the content of the ascii image.
func (m *Model) SetContent(content string) {
	m.Content = content
}

// SetImage sets the image of the ascii image.
func (m *Model) SetImage(img image.Image) {
	m.Image = img
}

// View returns a string representation of the ascii image.
func (m Model) View() string {
	return m.Content
}
