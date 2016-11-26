package main

import (
	"github.com/5k3105/Pipeline-Editor/genxml"
	"github.com/5k3105/Pipeline-Editor/gfxcanvas"
	"github.com/5k3105/Pipeline-Editor/gfxobjects/rectangle"
	"github.com/5k3105/Pipeline-Editor/graphpanel"
	"github.com/5k3105/Pipeline-Editor/janus"
	"os"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Application struct {
	*widgets.QApplication
	Window     *widgets.QMainWindow
	Canvas     *gfxcanvas.Canvas
	GraphPanel *graphpanel.GraphPanel
	Statusbar  *widgets.QStatusBar
}

func main() {
	var ap = &Application{}
	ap.QApplication = widgets.NewQApplication(len(os.Args), os.Args)

	ap.Window = widgets.NewQMainWindow(nil, 0)
	ap.Window.SetWindowTitle("Pipeline Editor")

	ap.Statusbar = widgets.NewQStatusBar(ap.Window)
	ap.Window.SetStatusBar(ap.Statusbar)

	ap.GraphPanel = graphpanel.NewGraphPanel(ap.Window, ap.CanvasReset, ap.LoadXML, ap.GenXML)
	ap.GraphPanel.SetMinimumWidth(320)

	ap.Canvas = gfxcanvas.NewCanvas(ap.Window, ap.ListNodes)
	ap.Canvas.Statusbar = ap.Statusbar

	ap.GraphPanel.SetJanusGraph(ap.Canvas.JanusGraph)

	ap.Window.AddDockWidget(core.Qt__LeftDockWidgetArea, ap.GraphPanel)
	ap.Window.SetCentralWidget(ap.Canvas.View)

	ap.Statusbar.ShowMessage(core.QCoreApplication_ApplicationDirPath(), 0)

	widgets.QApplication_SetStyle2("fusion")
	ap.Window.ShowMaximized()
	widgets.QApplication_Exec()
}

func (ap *Application) CanvasReset() {
	ap.Canvas = gfxcanvas.NewCanvas(ap.Window, ap.ListNodes)
	ap.Canvas.Statusbar = ap.Statusbar
	ap.Window.SetCentralWidget(ap.Canvas.View)
	ap.GraphPanel.SetJanusGraph(ap.Canvas.JanusGraph)

	//	point := core.NewQPointF3(300, 300)
	//	ap.Canvas.View.CenterOn(point)

	ap.GraphPanel.Clear()
}

func (ap *Application) GenXML(filename string) {
	genxml.GenXML(ap.GraphPanel, ap.Canvas.JanusGraph, filename)
}

func (ap *Application) LoadXML(filename string) {
	//ap.Canvas.Reset()
	ap.CanvasReset()
	genxml.LoadXML(filename, ap.Canvas, ap.GraphPanel)
	ap.ListNodes()

	// hack
	var strHeaders1 = []string{"Incoming Edge"}
	ap.GraphPanel.EdgeIncoming.SetHorizontalHeaderLabels(strHeaders1)
	ap.GraphPanel.EdgeIncoming.HorizontalHeader().SetStretchLastSection(true)

	var strHeaders2 = []string{"Outgoing Edge"}
	ap.GraphPanel.EdgeOutgoing.SetHorizontalHeaderLabels(strHeaders2)
	ap.GraphPanel.EdgeOutgoing.HorizontalHeader().SetStretchLastSection(true)

	ap.GraphPanel.NodeList.SetCurrentRow(0)

	//	point := core.NewQPointF3(300, 300)
	//	ap.Canvas.View.CenterOn(point)
}

func (ap *Application) TextUpdate() { // update text from tblargs to rect font

}

func (ap *Application) ListNodes() {
	ap.GraphPanel.NodeList.Clear()

	font := gui.NewQFont2("verdana", 13, 1, false)

	for _, f := range ap.Canvas.Figures.Values() {
		r := f.(*rectangle.Rectangle)

		ij, _ := ap.Canvas.JanusGraph.Nodes.Get(r.Node.Id)
		n := ij.(*janus.Node)

		var txt = n.ClassName
		r.Name = n.ClassName

		var li = widgets.NewQListWidgetItem2(txt, ap.GraphPanel.NodeList, 0)
		li.SetFont(font)
		li.SetData(int(core.Qt__UserRole), core.NewQVariant7(r.Node.Id))
	}

}
