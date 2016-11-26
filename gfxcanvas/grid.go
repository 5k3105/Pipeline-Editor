package gfxcanvas

import "github.com/therecipe/qt/gui"

type Grid struct {
	Color
	PenWidth           int
	X, Y               float64
	Vstep, Hstep, Step float64
}

func (c *Canvas) NewGrid(penwidth int, x, y, vstep, hstep, step float64, r, g, b, a int) *Grid {

	var inpcolor = Color{r, g, b, a}

	var grid = &Grid{
		Color:    inpcolor,
		PenWidth: penwidth,
		X:        x,
		Y:        y,
		Vstep:    vstep,
		Hstep:    hstep,
		Step:     step}

	var color = gui.NewQColor3(r, g, b, a)
	var pen = gui.NewQPen3(color)
	pen.SetWidth(penwidth)

	for x := grid.X; x <= hstep*step; x += step {
		c.Scene.AddLine2(x, 0, x, vstep*step, pen)
	}

	for y := grid.Y; y <= vstep*step; y += step {
		c.Scene.AddLine2(0, y, hstep*step, y, pen)
	}

	c.View.SetScene(c.Scene)
	c.View.Show()
	return grid
}

//	vstep := 24.0
//	hstep := 64.0
//	step := 5.0
