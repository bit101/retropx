package bitmap

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

type Bitmap struct {
	*image.RGBA
	Color color.RGBA
}

func NewBitmap(w, h int) *Bitmap {
	return &Bitmap{image.NewRGBA(image.Rect(0, 0, w, h)), color.RGBA{0, 0, 0, 255}}
}

func (b *Bitmap) Line(x0, y0, x1, y1 int) {
	if math.Abs(float64(y1-y0)) < math.Abs(float64(x1-x0)) {
		if x0 > x1 {
			x1, x0 = x0, x1
			y1, y0 = y0, y1
		}
		b.lineLow(x0, y0, x1, y1)
	} else {
		if y0 > y1 {
			x1, x0 = x0, x1
			y1, y0 = y0, y1
		}
		b.lineHigh(x0, y0, x1, y1)
	}
}

func (b *Bitmap) lineLow(x0, y0, x1, y1 int) {
	dx := x1 - x0
	dy := y1 - y0
	yi := 1
	if dy < 0 {
		yi = -1
		dy = -dy
	}
	d := 2*dy - dx
	y := y0
	for x := x0; x <= x1; x++ {
		b.Set(x, y, b.Color)
		if d > 0 {
			y += yi
			d -= (2 * dx)
		}
		d += (2 * dy)
	}
}

func (b *Bitmap) lineHigh(x0, y0, x1, y1 int) {
	dx := x1 - x0
	dy := y1 - y0
	xi := 1
	if dx < 0 {
		xi = -1
		dx = -dx
	}
	d := 2*dx - dy
	x := x0
	for y := y0; y <= y1; y++ {
		b.Set(x, y, b.Color)
		if d > 0 {
			x += xi
			d -= (2 * dy)
		}
		d += (2 * dx)
	}
}

func (b *Bitmap) Rect(x, y, w, h int) {
	b.Line(x, y, x+w, y)
	b.Line(x, y+h, x+w, y+h)
	b.Line(x, y, x, y+h)
	b.Line(x+w, y, x+w, y+h)
}

func (b *Bitmap) Square(x, y, size int) {
	b.Rect(x, y, size, size)
}

func (b *Bitmap) Circle(x, y, r int) {
	x0 := x + r
	y0 := y
	res := 10.0 / float64(r)
	for a := 0.0; a < math.Pi*2; a += res {
		x1 := x + int(math.Cos(a)*float64(r))
		y1 := y + int(math.Sin(a)*float64(r))
		b.Line(x0, y0, x1, y1)
		x0 = x1
		y0 = y1
	}
	b.Line(x0, y0, x+r, y)
}

func (b *Bitmap) Oval(x, y, w, h int) {
	rx := float64(w) / 2.0
	ry := float64(h) / 2.0
	x0 := x + int(rx)
	y0 := y
	res := 10.0 / math.Max(rx, ry)
	for a := 0.0; a < math.Pi*2; a += res {
		x1 := x + int(math.Cos(a)*rx)
		y1 := y + int(math.Sin(a)*ry)
		b.Line(x0, y0, x1, y1)
		x0 = x1
		y0 = y1
	}
	b.Line(x0, y0, x+int(rx), y)
}

func (b *Bitmap) Save(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("couldn't create image")
	}
	err = png.Encode(f, b)
	if err != nil {
		f.Close()
		log.Fatal("couldn't encode image")
	}
	f.Close()

}

func (bm *Bitmap) SetRGB(r, g, b uint8) {
	bm.Color = color.RGBA{r, g, b, 255}
}

func (bm *Bitmap) SetRGBA(r, g, b, a uint8) {
	bm.Color = color.RGBA{r, g, b, a}
}
