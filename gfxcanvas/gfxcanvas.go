package gfxcanvas

import (
	"github.com/5k3105/FM/gfxinterface"
	"github.com/5k3105/FM/gfxobjects/line"
	"github.com/5k3105/FM/gfxobjects/rectangle"
	"github.com/5k3105/FM/graph"

	"github.com/5k3105/FM/janus"
	"strconv"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

const Scalefactor = 10
const MvboxSize int = 100

var ListNodes func() // circular dependencies

var w, h float64 = 30 * Scalefactor, 10 * Scalefactor

type Canvas struct {
	Scene     *widgets.QGraphicsScene
	View      *widgets.QGraphicsView
	Statusbar *widgets.QStatusBar
	Graph     *graph.Graph
	*Grid
	Figures     *treemap.Map // gfxinterface.Figure
	drag        bool
	drawingline *line.Line
	move        bool
	movingrect  *rectangle.Rectangle
	JanusGraph  *janus.Graph
}

type Color struct{ R, G, B, A int }

func NewCanvas(parent *widgets.QMainWindow, lf func()) *Canvas {

	ListNodes = lf

	var canvas = &Canvas{
		Scene:      widgets.NewQGraphicsScene(parent),
		View:       widgets.NewQGraphicsView(parent),
		Figures:    treemap.NewWithIntComparator(),
		JanusGraph: janus.NewJanusGraph()}

	canvas.Grid = canvas.NewGrid(0, 100, 100, 64, 32, 5*Scalefactor, 160, 160, 160, 120)

	canvas.Scene.ConnectKeyPressEvent(canvas.keyPressEvent)
	canvas.Scene.ConnectMousePressEvent(canvas.mousePressEvent)
	canvas.Scene.ConnectWheelEvent(canvas.wheelEvent)
	canvas.Scene.ConnectMouseMoveEvent(canvas.mouseMoveEvent)

	canvas.Graph = graph.NewGraph()

	canvas.View.Scale(0.8, 0.8)
	canvas.View.Scale(0.8, 0.8)

	point := core.NewQPointF3(300, 300)
	canvas.View.CenterOn(point)

	canvas.View.SetViewportUpdateMode(0)

	// canvas.View.Scale(1, -1) // flips image

	return canvas
}
func (c *Canvas) Reset() {
	c.Scene.Clear()
	c.JanusGraph = janus.NewJanusGraph()
	c.Graph = graph.NewGraph()
	c.Figures.Clear()
	c.Grid = c.NewGrid(0, 100, 100, 64, 32, 5*Scalefactor, 160, 160, 160, 120)
	point := core.NewQPointF3(300, 300)
	c.View.CenterOn(point)
}

func (c *Canvas) mouseMoveEvent(e *widgets.QGraphicsSceneMouseEvent) {

	if c.move {
		var px, py = e.ScenePos().X(), e.ScenePos().Y()
		var x, y = float64((int(px) / 50) * 50), float64((int(py) / 50) * 50) // round to nearest fit

		c.movingrect.SetX(x - w + 50.0) // + float64(MvboxSize))
		c.movingrect.SetY(y)
		c.movingrect.PrepareGeometryChange()

	}

	if c.drag {
		var l = c.drawingline
		target := l.Target // Point
		target.SetX(e.ScenePos().X())
		target.SetY(e.ScenePos().Y())
		l.PrepareGeometryChange()

		c.Statusbar.ShowMessage(FloatToString(target.GetX())+" "+FloatToString(target.GetY()), 0)
	}

	c.Scene.MouseMoveEventDefault(e)
}

func (c *Canvas) wheelEvent(e *widgets.QGraphicsSceneWheelEvent) {
	if e.Modifiers() == core.Qt__ControlModifier {
		if e.Delta() > 0 {
			c.View.Scale(1.25, 1.25)
		} else {
			c.View.Scale(0.8, 0.8)
		}
	}
}

func (c *Canvas) mousePressEvent(e *widgets.QGraphicsSceneMouseEvent) {
	var px, py = e.ScenePos().X(), e.ScenePos().Y()

	switch e.Button() {
	case 1: // left button

		if c.move { // stop moving rectangle
			t, _ := c.JanusGraph.Nodes.Get(c.movingrect.GetNode().Id)
			tn := t.(*janus.Node)
			tn.X = c.movingrect.X
			tn.Y = c.movingrect.Y // save to jgraph

			c.movingrect.Over2 = false
			c.move = false
			c.movingrect = nil
			// do not DestroyQGraphicsItem() because this is a pointer to the real rectangle
			return
		}

		var x, y = float64((int(px) / 50) * 50), float64((int(py) / 50) * 50) // round to nearest fit
		//float64((int(px) / 5) * 5)
		var m = c.FindRectOverlap(x, y)
		if m == nil { // no node

			r := c.AddRectangle(x-w+50, y, w, h)
			c.Statusbar.ShowMessage(FloatToString(r.X)+" "+FloatToString(r.Y), 0)

		} else {

			var m = c.FindRectBeneath(px, py)
			if m != nil { // found node
				if c.drawingline == nil {

					r := m.(*rectangle.Rectangle)
					if r.FindRectBeneath(px, py) { // if over move box
						r.Over2 = true
						c.move = true
						c.movingrect = r

					} else {

						c.DrawLine(m, px, py)
						c.Statusbar.ShowMessage("drag start", 0)
					}
				} else {

					c.AddLine(c.drawingline.Source, m)
					c.Statusbar.ShowMessage("drag end", 0)

				}
			}
		}
	case 2: // right button

		if c.drag {
			c.CancelDraw()
		} else {
			var m = c.FindRectBeneath(px, py)
			if m != nil { // found node
				c.RemoveRectangle(m)
			}
		}

	}

	ListNodes() // callback

	c.Scene.MousePressEventDefault(e)

}

func (c *Canvas) RemoveRectangle(t gfxinterface.Figure) {

	var target = t.(*rectangle.Rectangle)
	c.RemoveEdges(target)

	if target.Scene().Pointer() != nil {
		target.DestroyQGraphicsItem()
		c.Figures.Remove(target.GetNode().Id) // remove from treemap

	}

	c.JanusGraph.RemoveNode(target.GetNode().Id)
}

func (c *Canvas) RemoveEdges(r *rectangle.Rectangle) {

	for _, f := range r.IncomingEdges.Values() { // **
		l := f.(*line.Line)
		if l.Scene().Pointer() != nil {
			l.DestroyQGraphicsPathItem()
		}
	}

	for _, f := range r.OutgoingEdges.Values() {
		l := f.(*line.Line)
		if l.Scene().Pointer() != nil {
			l.DestroyQGraphicsPathItem()
		}
	}

	r.IncomingEdges.Clear()
	r.OutgoingEdges.Clear() // remove from treemap

	nodeid := strconv.Itoa(r.GetNode().Id)

	// use iterators next

	for _, n := range c.JanusGraph.Nodes.Values() {
		for _, ne := range n.(*janus.Node).Edges.Values() {
			if ne.(*janus.Edge).String == nodeid {
				n.(*janus.Node).Edges.Remove(r.GetNode().Id)
			}
		}
		for _, ne := range n.(*janus.Node).AntiEdges.Values() {
			if ne.(*janus.Edge).String == nodeid {
				n.(*janus.Node).AntiEdges.Remove(r.GetNode().Id)
			}
		}
	}

}

func (c *Canvas) AddRectangle(x, y, w, h float64) *rectangle.Rectangle {

	i, r := rectangle.NewRectangle(x, y, w, h, c.Graph)
	c.Figures.Put(i, r)
	c.Scene.AddItem(r)

	c.JanusGraph.AddNode(i, x, y)

	return r
}

func (c *Canvas) AddRectangleFromFile(i int, x, y, w, h float64, language, classname, scriptfile string) *rectangle.Rectangle {

	i, r := rectangle.NewRectangleFromFile(i, x, y, w, h, c.Graph)
	c.Figures.Put(i, r)
	c.Scene.AddItem(r)

	c.JanusGraph.AddNodeFromFile(i, x, y, language, classname, scriptfile)

	return r
}

func (c *Canvas) AddLine(source gfxinterface.Figure, target gfxinterface.Figure) {

	l := line.AddLine(c.Graph, source, target)

	l.Source.AddEdgeOutgoing(l)
	l.Target.AddEdgeIncoming(l)

	c.Scene.AddItem(l)

	if c.drawingline.Scene().Pointer() != nil {
		c.drawingline.DestroyQGraphicsPathItem()
		c.drawingline = nil
		c.drag = false

	}

	t, _ := c.JanusGraph.Nodes.Get(target.GetNode().Id)
	tn := t.(*janus.Node)

	tn.AddEdge(source.GetNode().Id, "Outgoing Edge")

	s, _ := c.JanusGraph.Nodes.Get(source.GetNode().Id)
	sn := s.(*janus.Node)

	sn.AddEdge(target.GetNode().Id, "Incoming Edge")

}

func (c *Canvas) AddLineFromFile(source gfxinterface.Figure, target gfxinterface.Figure) {

	t, _ := c.JanusGraph.Nodes.Get(target.GetNode().Id)
	tn := t.(*janus.Node)

	tn.AddEdge(source.GetNode().Id, "Outgoing Edge")

	// s -> t
	s, _ := c.JanusGraph.Nodes.Get(source.GetNode().Id)
	sn := s.(*janus.Node)

	sn.AddEdge(target.GetNode().Id, "Incoming Edge")

	l := line.AddLine(c.Graph, source, target)

	c.Scene.AddItem(l)

	l.Source.AddEdgeOutgoing(l)
	l.Target.AddEdgeIncoming(l)

}

func (c *Canvas) DrawLine(source gfxinterface.Figure, tx, ty float64) {
	l := line.DrawLine(source, tx, ty)
	c.Scene.AddItem(l)
	c.drawingline = l
	c.drag = true
}

func (c *Canvas) keyPressEvent(e *gui.QKeyEvent) {

	if e.Modifiers() == core.Qt__ControlModifier {
		switch int32(e.Key()) {
		case int32(core.Qt__Key_Equal):
			c.View.Scale(1.25, 1.25)

		case int32(core.Qt__Key_Minus):
			c.View.Scale(0.8, 0.8)
		}
	}

	if e.Key() == int(core.Qt__Key_Escape) {
		if c.drag {
			c.CancelDraw()
		}
	}
}

func (c *Canvas) CancelDraw() {
	if c.drawingline.Scene().Pointer() != nil {
		c.drawingline.DestroyQGraphicsPathItem()
		c.drawingline = nil
		c.drag = false

	}
}

func (c *Canvas) ClearScene() {
	c.Scene.Clear()
	c.View.SetScene(c.Scene)
	c.View.Show()

}

func (c *Canvas) ShowPic(filepath, filetype string) {

	ir := gui.NewQImageReader3(filepath, filetype)
	img := ir.Read()

	pix := gui.QPixmap_FromImage(img, 0)

	c.Scene.Clear()
	c.Scene.AddPixmap(pix)

	c.View.SetScene(c.Scene)
	c.View.Show()

}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 1, 64)
}
