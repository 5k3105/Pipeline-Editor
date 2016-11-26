package line

import (
	"github.com/5k3105/Pipeline-Editor/gfxinterface"
	"github.com/5k3105/Pipeline-Editor/gfxobjects/point"
	"github.com/5k3105/Pipeline-Editor/graph"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const penwidth int = 5

type Line struct {
	*widgets.QGraphicsPathItem
	Edge    *graph.Edge // ?
	Type    string
	Source  gfxinterface.Figure
	Target  gfxinterface.Figure
	SX, SY  float64 // source
	TX, TY  float64 // target
	fgcolor Color
	over    bool
}

type Color struct{ R, G, B, A int }

func AddLine(g *graph.Graph, source gfxinterface.Figure, target gfxinterface.Figure) *Line {

	var s, t = source.GetNode(), target.GetNode()

	var l = &Line{
		QGraphicsPathItem: widgets.NewQGraphicsPathItem(nil),
		Type:              "Line",
		Edge:              g.AddEdge(s, t),
		Source:            source,
		Target:            target,
		SX:                source.GetX(),
		SY:                source.GetY(),
		TX:                target.GetX(),
		TY:                target.GetY()}

	l.ConnectBoundingRect(l.boundingRect)
	l.ConnectPaint(l.paint)

	return l
}

//func AddLineFromFile(g *graph.Graph, source gfxinterface.Figure, target gfxinterface.Figure) *Line {

//	var s, t = source.GetNode(), target.GetNode()

//	var l = &Line{
//		QGraphicsPathItem: widgets.NewQGraphicsPathItem(nil),
//		Type:              "Line",
//		Edge:              g.AddEdge(s, t),
//		Source:            source,
//		Target:            target,
//		SX:                source.GetX(),
//		SY:                source.GetY(),
//		TX:                target.GetX(),
//		TY:                target.GetY()}

//	l.ConnectBoundingRect(l.boundingRect)
//	l.ConnectPaint(l.paint)

//	return l
//}

func (l *Line) GetEdge() *graph.Edge { return l.Edge }

func (l *Line) paint(p *gui.QPainter, o *widgets.QStyleOptionGraphicsItem, widget *widgets.QWidget) {
	var color = gui.NewQColor3(0, 0, 0, 255)
	var brush = gui.NewQBrush3(color, 0) // solid = 1, nobrush = 0
	var pen = gui.NewQPen3(color)
	pen.SetWidth(penwidth)

	p.SetRenderHint(1, true) // Antiailiasing
	p.SetPen(pen)
	p.SetBrush(brush)

	path := gui.NewQPainterPath()

	source := l.Source
	target := l.Target

	var sx, sy, tx, ty, w = source.GetX(), source.GetY(), target.GetX(), target.GetY(), target.GetW()
	sx, sy = sx+(w/2.0), sy+100.0 // bottom center is output node
	tx, ty = tx+(w/2.0), ty       // top center is input node
	path.MoveTo2(sx, sy)
	sy = sy + 50.0 // offset by single step
	path.LineTo2(sx, sy)

	//var A, B, C int

	// target straight down
	if sy < ty && sx == tx {
		path.LineTo2(tx, ty)
		p.DrawPath(path)
		drawArrow(p, tx, ty)
		return
	}

	// target straight up
	if sy > ty && sx == tx {
		var offset = (w/2 + 50.0)
		sx = sx + offset
		path.LineTo2(sx, sy) // A
		ty = ty - 50.0
		path.LineTo2(sx, ty) // B
		path.LineTo2(tx, ty) // C

		ty = ty + 50.0
		path.LineTo2(tx, ty)
		p.DrawPath(path)
		drawArrow(p, tx, ty)
		return
	}

	// target bottom left
	if sx > tx && sy < ty {
		path.LineTo2(tx, sy)
		path.LineTo2(tx, ty)
		p.DrawPath(path)
		drawArrow(p, tx, ty)
		return
	}

	// target bottom right
	if sx < tx && sy < ty {
		path.LineTo2(tx, sy)
		path.LineTo2(tx, ty)
		p.DrawPath(path)
		drawArrow(p, tx, ty)
		return
	}

	// target top right
	if sx < tx && sy > ty {
		var offset = (w/2 + 50.0)
		if sx+offset <= tx-offset { // between
			sx = sx + offset
			path.LineTo2(sx, sy) // A

			ty = ty - 50.0
			path.LineTo2(sx, ty) // B
			path.LineTo2(tx, ty) // C

			ty = ty + 50.0
			path.LineTo2(tx, ty)
			p.DrawPath(path)
			drawArrow(p, tx, ty)
			return
		} else { // around
			tx = tx + offset
			path.LineTo2(tx, sy) // A

			ty = ty - 50.0
			path.LineTo2(tx, ty) // B
			tx = tx - offset
			path.LineTo2(tx, ty) // C

			ty = ty + 50.0
			path.LineTo2(tx, ty)
			p.DrawPath(path)
			drawArrow(p, tx, ty)
			return
		}
	}

	// target top left
	if sx > tx && sy > ty {
		var offset = (w/2 + 50.0)
		if sx-offset >= tx+offset { // between
			sx = sx - offset
			path.LineTo2(sx, sy) // A

			ty = ty - 50.0
			path.LineTo2(sx, ty) // B
			path.LineTo2(tx, ty) // C

			ty = ty + 50.0
			path.LineTo2(tx, ty)
			p.DrawPath(path)
			drawArrow(p, tx, ty)
			return
		} else { // around
			tx = tx - offset
			path.LineTo2(tx, sy) // A

			ty = ty - 50.0
			path.LineTo2(tx, ty) // B
			tx = tx + offset
			path.LineTo2(tx, ty) // C

			ty = ty + 50.0
			path.LineTo2(tx, ty)
			p.DrawPath(path)
			drawArrow(p, tx, ty)
			return
		}
	}

}

func drawArrow(p *gui.QPainter, x, y float64) {
	var color = gui.NewQColor3(0, 0, 0, 255)
	var brush = gui.NewQBrush3(color, 1) // solid = 1, nobrush = 0
	var pen = gui.NewQPen3(color)
	pen.SetWidth(0)

	p.SetRenderHint(1, true) // Antiailiasing
	p.SetPen(pen)
	p.SetBrush(brush)

	path := gui.NewQPainterPath()

	path.MoveTo2(x-10, y-20)
	path.LineTo2(x+10, y-20)
	path.LineTo2(x, y)
	path.LineTo2(x-10, y-20)

	p.DrawPath(path)
}

func (l *Line) boundingRect() *core.QRectF {
	source := l.Source
	target := l.Target
	return core.NewQRectF4(source.GetX(), source.GetY(), target.GetX(), target.GetY())
}

func (l *Line) boundingRectDL() *core.QRectF {
	source := l.Source // Rectangle
	target := l.Target // Point
	return core.NewQRectF4(source.GetX(), source.GetY(), target.GetX(), target.GetY())
}

// draw temporary line
func DrawLine(source gfxinterface.Figure, tx, ty float64) *Line {

	var target = &point.Point{X: tx, Y: ty}

	l := &Line{
		QGraphicsPathItem: widgets.NewQGraphicsPathItem(nil),
		Type:              "Line",
		Source:            source,
		Target:            target}

	l.ConnectBoundingRect(l.boundingRectDL)
	l.ConnectPaint(l.draw)

	return l
}

func (l *Line) draw(p *gui.QPainter, o *widgets.QStyleOptionGraphicsItem, widget *widgets.QWidget) {
	var color = gui.NewQColor3(0, 0, 0, 255)
	var brush = gui.NewQBrush3(color, 0)
	var pen = gui.NewQPen3(color)
	pen.SetWidth(penwidth)

	p.SetRenderHint(1, true) // Antiailiasing
	p.SetPen(pen)
	p.SetBrush(brush)

	path := gui.NewQPainterPath()
	// center line start, extend to current cursor location
	source := l.Source                                               // Rectangle
	target := l.Target                                               // Point
	path.MoveTo2(source.GetX()+source.GetW()/2, source.GetY()+100/2) // GetH() ??
	path.LineTo2(target.GetX(), target.GetY())

	p.DrawPath(path)
}

//func AutoNumber() int64 {
//	autonumber += 1
//	return autonumber
//}

//------------------------------------------------------------------------------

//type Line struct {
//	*widgets.QGraphicsPathItem
//	Color
//	Id     int
//	Source Shape
//	Target Shape
//	//	SId, TId int
//	//	SX, SY   float64 // source
//	//	TX, TY   float64 // target
//	over bool
//}

//		SX:                sx,
//		SY:                sy,
//		TX:                tx,
//		TY:                ty}

// paint:
//	l.SetPath(path)
//	l.Scene.AddPath(path, pen, brush)
