package main

import "retropx/bitmap"

func main() {
	b := bitmap.NewBitmap(600, 600)

	// b.Line(500, 500, 100, 100)
	// b.Circle(300, 300, 250)
	// b.Oval(400, 400, 150, 100)
	// b.Rect(10, 10, 100, 200)
	// b.Square(20, 20, 100)
	// b.FillRect(10, 10, 100, 100)
	// b.FillCircle(200, 200, 50)
	// b.BezierCurve(100, 100, 600, 0, 600, 600, 0, 300)
	// b.QuadraticCurve(100, 100, 600, 0, 600, 600)
	b.Arc(300, 300, 100, 0, 270)

	b.Save("out.png")
}
