package rectangle

import (
	"github.com/5k3105/Pipeline-Editor/gfxinterface"
	"github.com/5k3105/Pipeline-Editor/graph"
	"strconv"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const MvboxSize float64 = 100
const penwidth int = 5

type Rectangle struct {
	*widgets.QGraphicsItem
	Type          string
	Name          string
	Node          *graph.Node
	IncomingEdges *treemap.Map //[]gfxinterface.Link
	OutgoingEdges *treemap.Map //[]gfxinterface.Link
	X, Y          float64
	W, H          float64
	over          bool
	Over2         bool
	fgcolor       Color
	bgcolor       Color
}

type Color struct{ R, G, B, A int }

func (r *Rectangle) GetX() float64        { return r.X }
func (r *Rectangle) SetX(x float64)       { r.X = x }
func (r *Rectangle) GetY() float64        { return r.Y }
func (r *Rectangle) SetY(y float64)       { r.Y = y }
func (r *Rectangle) GetW() float64        { return r.W }
func (r *Rectangle) GetNode() *graph.Node { return r.Node }

func (r *Rectangle) AddEdgeIncoming(l gfxinterface.Link) {
	e := l.GetEdge()
	r.IncomingEdges.Put(e.Id, l)
}

func (r *Rectangle) AddEdgeOutgoing(l gfxinterface.Link) {
	e := l.GetEdge()
	r.OutgoingEdges.Put(e.Id, l)
}

func NewRectangle(x, y, w, h float64, g *graph.Graph) (int, *Rectangle) {

	r := &Rectangle{
		QGraphicsItem: widgets.NewQGraphicsItem(nil),
		Type:          "Rectangle",
		Node:          g.AddNode(),
		IncomingEdges: treemap.NewWithIntComparator(),
		OutgoingEdges: treemap.NewWithIntComparator(),
		X:             x,
		Y:             y,
		W:             w,
		H:             h,
		fgcolor:       Color{0, 0, 0, 255},
		bgcolor:       Color{0, 0, 0, 100},
	}

	r.SetAcceptHoverEvents(true)
	r.ConnectHoverEnterEvent(r.hoverEnterEvent)
	r.ConnectHoverLeaveEvent(r.hoverLeaveEvent)

	r.ConnectBoundingRect(r.boundingRect)
	r.ConnectPaint(r.paint)

	return r.Node.Id, r
}

func NewRectangleFromFile(i int, x, y, w, h float64, g *graph.Graph) (int, *Rectangle) {

	r := &Rectangle{
		QGraphicsItem: widgets.NewQGraphicsItem(nil),
		Type:          "Rectangle",
		Node:          g.AddNodeFromFile(i),
		IncomingEdges: treemap.NewWithIntComparator(),
		OutgoingEdges: treemap.NewWithIntComparator(),
		X:             x,
		Y:             y,
		W:             w,
		H:             h,
		fgcolor:       Color{0, 0, 0, 255},
		bgcolor:       Color{0, 0, 0, 100},
	}

	r.SetAcceptHoverEvents(true)
	r.ConnectHoverEnterEvent(r.hoverEnterEvent)
	r.ConnectHoverLeaveEvent(r.hoverLeaveEvent)

	r.ConnectBoundingRect(r.boundingRect)
	r.ConnectPaint(r.paint)

	return r.Node.Id, r
}

func (r *Rectangle) paint(p *gui.QPainter, o *widgets.QStyleOptionGraphicsItem, w *widgets.QWidget) {

	var color = gui.NewQColor3(r.fgcolor.R, r.fgcolor.G, r.fgcolor.B, r.fgcolor.A) // r,g,b,a

	var pen = gui.NewQPen3(color)
	pen.SetWidth(penwidth)

	if r.Over2 {
		pen.SetStyle(0)
	}

	p.SetRenderHint(1, true) // Antiailiasing
	var path = gui.NewQPainterPath()
	path.AddRoundedRect2(r.X, r.Y, r.W, r.H, 1, 1, 0) // Qt::AbsoluteSize
	p.SetPen(pen)
	p.DrawPath(path)

	if r.over {
		color = gui.NewQColor3(r.bgcolor.R, r.bgcolor.G, r.bgcolor.B, r.bgcolor.A) // r,g,b,a
		var brush = gui.NewQBrush3(color, 1)
		p.FillPath(path, brush)

		color = gui.NewQColor3(0, 0, 0, 60) // r,g,b,a
		var pen2 = gui.NewQPen2(0)          // no pen

		var brush2 = gui.NewQBrush3(color, 1)
		var path2 = gui.NewQPainterPath()

		p.SetPen(pen2)
		path2.AddRoundedRect2(r.X+r.W-MvboxSize, r.Y, MvboxSize, MvboxSize, 1, 1, 0) // moveme box // (r.X+r.W-5, r.Y, 5, 5, 1, 1, 0)
		p.FillPath(path2, brush2)
		p.DrawPath(path2)
	}

	// if selected {}

	var font = gui.NewQFont2("verdana", 20, 1, false)
	p.SetFont(font)

	var qpf = core.NewQPointF3(r.X+10.0, r.Y+r.H-10.0)

	if r.Name == "" {
		p.DrawText(qpf, "Node "+strconv.Itoa(r.Node.Id))
	} else {
		p.DrawText(qpf, r.Name) //+" "+strconv.Itoa(r.Node.Id))
	}
	//p.DrawPath(path)
}

func (r *Rectangle) boundingRect() *core.QRectF {
	return core.NewQRectF4(r.X, r.Y, r.W, r.H)
}

func (r *Rectangle) hoverEnterEvent(e *widgets.QGraphicsSceneHoverEvent) {
	r.over = true
	r.Update(core.NewQRectF4(r.X, r.Y, r.W, r.H))
}

func (r *Rectangle) hoverLeaveEvent(e *widgets.QGraphicsSceneHoverEvent) {
	r.over = false
	r.Update(core.NewQRectF4(r.X, r.Y, r.W, r.H))
}

func (r *Rectangle) FindRectBeneath(x, y float64) bool {
	var iX, iY, iW = r.X, r.Y, r.W
	//if x >= iX+iW-5 && x <= iX+iW && y >= iY && y <= iY+5 {
	if x >= iX+iW-MvboxSize && x <= iX+iW && y >= iY && y <= iY+MvboxSize {
		return true
	} else {
		return false
	}

}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}
