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

func (b *Bitmap) Plot(x, y int) {
	b.Set(x, y, b.Color)
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
		b.Plot(x, y)
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
		b.Plot(x, y)
		if d > 0 {
			x += xi
			d -= (2 * dy)
		}
		d += (2 * dx)
	}
}

func (b *Bitmap) Rect(x, y, w, h int) {
	for xx := x; xx < x+w; xx++ {
		b.Plot(xx, y)
		b.Plot(xx, y+h)
	}
	for yy := y; yy < y+h; yy++ {
		b.Plot(x, yy)
		b.Plot(x+w, yy)
	}
}

func (b *Bitmap) FillRect(x, y, w, h int) {
	for yy := y; yy < y+h; yy++ {
		for xx := x; xx < x+w; xx++ {
			b.Plot(xx, yy)
		}
	}
}

func (b *Bitmap) Square(x, y, size int) {
	b.Rect(x, y, size, size)
}

func (b *Bitmap) FillSquare(x, y, size int) {
	b.FillRect(x, y, size, size)
}

func (b *Bitmap) Circle(x, y, radius int) {
	xx := radius - 1
	yy := 0
	dx := 1
	dy := 1
	err := dx - (radius << 1)

	for xx >= yy {
		b.Plot(x+xx, y+yy)
		b.Plot(x+yy, y+xx)
		b.Plot(x-yy, y+xx)
		b.Plot(x-xx, y+yy)
		b.Plot(x-xx, y-yy)
		b.Plot(x-yy, y-xx)
		b.Plot(x+yy, y-xx)
		b.Plot(x+xx, y-yy)

		if err <= 0 {
			yy++
			err += dy
			dy += 2
		}

		if err > 0 {
			xx--
			dx += 2
			err += dx - (radius << 1)
		}
	}
}

func (b *Bitmap) FillCircle(x, y, radius int) {
	xx := radius - 1
	yy := 0
	dx := 1
	dy := 1
	err := dx - (radius << 1)

	for xx >= yy {
		for xa := x - xx; xa <= x+xx; xa++ {
			b.Plot(xa, y+yy)
			b.Plot(xa, y-yy)
		}
		for xa := x - yy; xa <= x+yy; xa++ {
			b.Plot(xa, y+xx)
			b.Plot(xa, y-xx)
		}
		if err <= 0 {
			yy++
			err += dy
			dy += 2
		}

		if err > 0 {
			xx--
			dx += 2
			err += dx - (radius << 1)
		}
	}
}

func (b *Bitmap) Arc(x, y, r, start, end int) {
	startRad := float64(start) * math.Pi / 180.0
	endRad := float64(end) * math.Pi / 180.0
	fr := float64(r)
	x0 := x + int(math.Cos(startRad)*fr)
	y0 := y + int(math.Sin(startRad)*fr)
	res := 10.0 / fr
	for a := startRad; a < endRad; a += res {
		x1 := x + int(math.Cos(a)*fr)
		y1 := y + int(math.Sin(a)*fr)
		b.Line(x0, y0, x1, y1)
		x0 = x1
		y0 = y1
	}
	b.Line(x0, y0, x+int(math.Cos(endRad)*fr), y+int(math.Sin(endRad)*fr))
}

func (b *Bitmap) Oval(x, y, rx, ry int) {
	frx := float64(rx)
	fry := float64(ry)
	x0 := x + rx
	y0 := y
	res := 10.0 / math.Max(frx, fry)
	for a := 0.0; a < math.Pi*2; a += res {
		x1 := x + int(math.Cos(a)*frx)
		y1 := y + int(math.Sin(a)*fry)
		b.Line(x0, y0, x1, y1)
		x0 = x1
		y0 = y1
	}
	b.Line(x0, y0, x+int(rx), y)
}

// func (b *Bitmap) FillOval(x, y, rx, ry int) {
// 	for xx := x - rx; xx <= x+rx; xx++ {
// 		for yy := y - ry; yy <= y+ry; yy++ {
// 			if math.Hypot(float64(xx-x), float64(yy-y)) < float64(r) {
// 				b.Plot(xx, yy)
// 			}
// 		}
// 	}
// }

func (b *Bitmap) BezierCurve(x0, y0, x1, y1, x2, y2, x3, y3 int) {
	fx0 := float64(x0)
	fy0 := float64(y0)
	fx1 := float64(x1)
	fy1 := float64(y1)
	fx2 := float64(x2)
	fy2 := float64(y2)
	fx3 := float64(x3)
	fy3 := float64(y3)
	x := x0
	y := y0
	for t := 0.0; t < 1.0; t += 0.01 {
		oneMinusT := 1.0 - t
		m0 := oneMinusT * oneMinusT * oneMinusT
		m1 := 3.0 * oneMinusT * oneMinusT * t
		m2 := 3.0 * oneMinusT * t * t
		m3 := t * t * t
		xt := int(m0*fx0 + m1*fx1 + m2*fx2 + m3*fx3)
		yt := int(m0*fy0 + m1*fy1 + m2*fy2 + m3*fy3)
		b.Line(x, y, xt, yt)
		x = xt
		y = yt
	}
}

func (b *Bitmap) QuadraticCurve(x0, y0, x1, y1, x2, y2 int) {
	fx0 := float64(x0)
	fy0 := float64(y0)
	fx1 := float64(x1)
	fy1 := float64(y1)
	fx2 := float64(x2)
	fy2 := float64(y2)
	x := x0
	y := y0
	for t := 0.0; t < 1.0; t += 0.01 {
		oneMinusT := 1.0 - t
		m0 := oneMinusT * oneMinusT
		m1 := 2.0 * oneMinusT * t
		m2 := t * t
		xt := int(m0*fx0 + m1*fx1 + m2*fx2)
		yt := int(m0*fy0 + m1*fy1 + m2*fy2)
		b.Line(x, y, xt, yt)
		x = xt
		y = yt
	}
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
