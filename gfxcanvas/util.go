package gfxcanvas

import (
	"github.com/5k3105/FM/gfxinterface"
	"github.com/5k3105/FM/gfxobjects/rectangle"
)

func (c *Canvas) FindRectOverlap(x, y float64) gfxinterface.Figure {
	// overlap of 2 rectangles

	// detect right and below
	//	L2x, L2y := x, y
	//	R2x, R2y := x+300, y+100

	// detect left and below
	L2x, L2y := x-300+50, y
	R2x, R2y := x, y+100

	for _, f := range c.Figures.Values() {
		i := f.(*rectangle.Rectangle)

		L1x, L1y := i.X, i.Y
		R1x, R1y := i.X+i.W, i.Y+i.H

		if L2x < R1x && R2x > L1x && L2y < R1y && R2y > L1y {
			return i
		}
	}

	return nil
}

func (c *Canvas) FindRectBeneath(x, y float64) gfxinterface.Figure {
	// point within another rect area
	for _, f := range c.Figures.Values() {
		i := f.(*rectangle.Rectangle)

		if x >= i.X && x <= i.X+i.W && y >= i.Y && y <= i.Y+i.H {
			return i
		}
	}

	return nil
}

//func (c *Canvas) FindRectBeneath2(x, y, iX, iY float64) bool {

//	if x >= iX+15 && x <= iX+20 && y >= iY && y <= iY+5 {
//		return true
//	} else {
//		return false
//	}

//}
