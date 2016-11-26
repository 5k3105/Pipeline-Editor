package graphpanel

import (
	"github.com/5k3105/FM/janus"

	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var (
	LoadXML     func(s string)
	GenXML      func(s string)
	CanvasReset func()
	graphpanel  *GraphPanel
	janusgraph  *janus.Graph
)

type GraphPanel struct {
	*widgets.QDockWidget
	FileList     Filelist
	FileName     Lineedit
	GraphName    Lineedit // QLineEdit
	JobID        Lineedit
	DatasetPath  Lineedit
	NodeID       Lineedit
	Language     Lineedit
	ScriptFile   Lineedit
	ClassName    Lineedit
	NodeList     Nodelist // QListWidget
	EdgeIncoming *widgets.QTableWidget
	EdgeOutgoing *widgets.QTableWidget
	TblArgs      *widgets.QTableView
	ValueList    Valuelist
	GenerateXML  *widgets.QPushButton
}

func NewGraphPanel(parent *widgets.QMainWindow, reset func(), lxml func(s string), genxml func(s string)) *GraphPanel {

	LoadXML = lxml
	GenXML = genxml
	CanvasReset = reset

	graphpanel = &GraphPanel{}
	graphpanel.QDockWidget = widgets.NewQDockWidget("Graph Panel", parent, 0)

	var w = widgets.NewQWidget(parent, 0)
	w.SetLayout(graphpanel.ui())
	graphpanel.SetWidget(w)

	return graphpanel
}

func (gp *GraphPanel) Clear() {
	gp.FileName.Clear()
	gp.GraphName.Clear()
	gp.JobID.Clear()
	gp.DatasetPath.Clear()
	gp.NodeID.Clear()
	gp.Language.Clear()
	gp.ScriptFile.Clear()
	gp.ClassName.Clear()
	gp.NodeList.Clear()
	gp.EdgeIncoming.Clear()
	gp.EdgeOutgoing.Clear()
	gp.Tblargs_Clear()
	//gp.TblArgs.Clear()
	gp.ValueList.Clear()

}

func (gp *GraphPanel) ui() *widgets.QVBoxLayout {

	var vlayout0 = widgets.NewQVBoxLayout()

	var lblFileset = widgets.NewQLabel2("File Set: ", nil, 0)

	gp.FileList = Filelist{widgets.NewQListWidget(nil)}
	gp.FileList.SetFont(gui.NewQFont2("verdana", 13, 1, false))
	gp.FileList.SetAcceptDrops(true)
	//gp.FileList.SetMaximumHeight(200)

	gp.FileList.ConnectDragEnterEvent(gp.FileList.DragEnterEvent_)
	gp.FileList.ConnectDragMoveEvent(gp.FileList.DragMoveEvent_)
	gp.FileList.ConnectDropEvent(gp.FileList.DropEvent_)
	gp.FileList.ConnectKeyPressEvent(gp.FileList.KeyPressEvent)

	gp.FileList.ConnectCurrentItemChanged(gp.FileList.CurrentItemChanged)

	vlayout0.AddWidget(lblFileset, 0, 0)
	vlayout0.AddWidget(gp.FileList, 0, 0)

	var hlayoutLine0 = widgets.NewQHBoxLayout()

	var lblFileName = widgets.NewQLabel2(" File Name: ", nil, 0)
	lblFileName.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	gp.FileName = Lineedit{widgets.NewQLineEdit(nil)}
	gp.FileName.SetObjectName("FileName")
	gp.FileName.ConnectEditingFinished(gp.FileName.EditingFinished)

	gp.GenerateXML = widgets.NewQPushButton(nil)
	gp.GenerateXML.SetText("Save")
	gp.GenerateXML.SetMaximumWidth(150)
	gp.GenerateXML.ConnectClicked(func(_ bool) { GenXML(gp.FileName.Text()) })

	hlayoutLine0.AddWidget(lblFileName, 0, 0)
	hlayoutLine0.AddWidget(gp.FileName, 0, 0)
	hlayoutLine0.AddWidget(gp.GenerateXML, 0, 0)

	var hlayoutLine1 = gp.line1() // graph
	var hlayoutLine2 = gp.line2() // node
	var vlayoutLine3 = gp.line3() // arg

	var vlayout = widgets.NewQVBoxLayout()

	vlayout.AddLayout(vlayout0, 0)
	vlayout.AddLayout(hlayoutLine0, 0)

	vlayout.AddLayout(hlayoutLine1, 0)
	vlayout.AddLayout(hlayoutLine2, 0)
	vlayout.AddLayout(vlayoutLine3, 0)

	return vlayout

}

func (gp *GraphPanel) SetJanusGraph(j *janus.Graph) {
	janusgraph = j
}

func (gp *GraphPanel) line1() *widgets.QHBoxLayout {
	var gbGraphProperties = widgets.NewQGroupBox2("Graph Properties", nil)

	var vlayout = widgets.NewQVBoxLayout()

	var hlayoutLine1 = widgets.NewQHBoxLayout() // graphName, jobID
	var hlayoutLine2 = widgets.NewQHBoxLayout() // datasetPath

	var lblGraphName = widgets.NewQLabel2(" Graph Name: ", nil, 0)
	lblGraphName.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	gp.GraphName = Lineedit{widgets.NewQLineEdit(nil)}
	gp.GraphName.SetObjectName("GraphName")
	gp.GraphName.ConnectEditingFinished(gp.GraphName.EditingFinished)

	var lblJobID = widgets.NewQLabel2(" Job ID: ", nil, 0)
	lblJobID.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	gp.JobID = Lineedit{widgets.NewQLineEdit(nil)}
	gp.JobID.SetObjectName("JobID")
	gp.JobID.ConnectEditingFinished(gp.JobID.EditingFinished)

	var btnNewGraph = widgets.NewQPushButton(nil)
	btnNewGraph.SetText("New Graph")
	//btnScriptFile.SetMaximumWidth(120)
	btnNewGraph.ConnectClicked(func(_ bool) { CanvasReset() })

	var lblDatasetPath = widgets.NewQLabel2(" Dataset Path: ", nil, 0)
	lblDatasetPath.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	gp.DatasetPath = Lineedit{widgets.NewQLineEdit(nil)}
	gp.DatasetPath.SetObjectName("DatasetPath")
	gp.DatasetPath.ConnectEditingFinished(gp.DatasetPath.EditingFinished)

	var vlayoutAlign = widgets.NewQVBoxLayout()
	var hlayoutAlign = widgets.NewQHBoxLayout()

	vlayoutAlign.AddWidget(lblGraphName, 0, 0)
	hlayoutLine1.AddWidget(gp.GraphName, 0, 0)
	hlayoutLine1.AddWidget(lblJobID, 0, 0)
	hlayoutLine1.AddWidget(gp.JobID, 0, 0)
	hlayoutLine1.AddWidget(btnNewGraph, 0, 0)

	vlayoutAlign.AddWidget(lblDatasetPath, 0, 0)
	hlayoutLine2.AddWidget(gp.DatasetPath, 0, 0)

	vlayout.AddLayout(hlayoutLine1, 0)
	vlayout.AddLayout(hlayoutLine2, 0)

	hlayoutAlign.AddLayout(vlayoutAlign, 0)
	hlayoutAlign.AddLayout(vlayout, 0)

	gbGraphProperties.SetLayout(hlayoutAlign)

	var hlayout = widgets.NewQHBoxLayout()
	hlayout.AddWidget(gbGraphProperties, 0, 0)

	return hlayout

}

func (gp *GraphPanel) line2() *widgets.QHBoxLayout {

	var hlayoutLine1 = widgets.NewQHBoxLayout() // nodeList
	var vlayoutLine2 = widgets.NewQVBoxLayout() // Node + Edges
	var hlayoutLine3 = widgets.NewQHBoxLayout() // List + Node/Edge

	var hlayoutLine21Node = widgets.NewQHBoxLayout() // Node
	var hlayoutLine22Edge = widgets.NewQHBoxLayout() // Edges

	// Node List
	gp.NodeList = Nodelist{widgets.NewQListWidget(nil)}
	gp.NodeList.ConnectCurrentItemChanged(gp.NodeList.CurrentItemChanged_)

	hlayoutLine1.AddWidget(gp.NodeList, 0, 0)

	//----------
	var gbNodeProperties = widgets.NewQGroupBox2("Node Properties", nil)

	var vlayout = widgets.NewQVBoxLayout()

	var hlayoutNpLine1 = widgets.NewQHBoxLayout() // nodeID, language
	var hlayoutNpLine2 = widgets.NewQHBoxLayout() // scriptfile
	var hlayoutNpLine3 = widgets.NewQHBoxLayout() // classname

	var lblNodeID = widgets.NewQLabel2(" Node ID: ", nil, 0)
	lblNodeID.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	gp.NodeID = Lineedit{widgets.NewQLineEdit(nil)}
	gp.NodeID.SetObjectName("NodeID")
	gp.NodeID.ConnectEditingFinished(gp.NodeID.EditingFinished)

	var lblLanguage = widgets.NewQLabel2(" Language: ", nil, 0)
	lblLanguage.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	gp.Language = Lineedit{widgets.NewQLineEdit(nil)}
	gp.Language.SetObjectName("Language")
	gp.Language.ConnectEditingFinished(gp.Language.EditingFinished)

	var lblScriptFile = widgets.NewQLabel2(" Script File: ", nil, 0)
	lblScriptFile.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	gp.ScriptFile = Lineedit{widgets.NewQLineEdit(nil)}
	gp.ScriptFile.SetObjectName("ScriptFile")
	gp.ScriptFile.ConnectEditingFinished(gp.ScriptFile.EditingFinished)

	var lblClassName = widgets.NewQLabel2(" Class Name: ", nil, 0)
	lblClassName.SetSizePolicy2(widgets.QSizePolicy__Fixed, widgets.QSizePolicy__Fixed)

	gp.ClassName = Lineedit{widgets.NewQLineEdit(nil)}
	gp.ClassName.SetObjectName("ClassName")
	gp.ClassName.ConnectEditingFinished(gp.ClassName.EditingFinished)

	var btnScriptFile = widgets.NewQPushButton(nil)
	btnScriptFile.SetText("+")
	btnScriptFile.SetMaximumWidth(35)
	btnScriptFile.ConnectClicked(func(_ bool) {
		if btnScriptFile.Text() == "+" {
			gp.EdgeIncoming.SetVisible(true)
			gp.EdgeOutgoing.SetVisible(true)
			btnScriptFile.SetText("-")
		} else {
			gp.EdgeOutgoing.SetVisible(false)
			gp.EdgeIncoming.SetVisible(false)
			btnScriptFile.SetText("+")
		}

	})

	var vlayoutAlign = widgets.NewQVBoxLayout()
	var hlayoutAlign = widgets.NewQHBoxLayout()

	vlayoutAlign.AddWidget(lblNodeID, 0, 0)
	hlayoutNpLine1.AddWidget(gp.NodeID, 0, 0)
	hlayoutNpLine1.AddWidget(lblLanguage, 0, 0)
	hlayoutNpLine1.AddWidget(gp.Language, 0, 0)

	vlayoutAlign.AddWidget(lblScriptFile, 0, 0)
	hlayoutNpLine2.AddWidget(gp.ScriptFile, 0, 0)
	hlayoutNpLine2.AddWidget(btnScriptFile, 0, 0)

	vlayoutAlign.AddWidget(lblClassName, 0, 0)
	hlayoutNpLine3.AddWidget(gp.ClassName, 0, 0)

	vlayout.AddLayout(hlayoutNpLine1, 0)
	vlayout.AddLayout(hlayoutNpLine2, 0)
	vlayout.AddLayout(hlayoutNpLine3, 0)

	hlayoutAlign.AddLayout(vlayoutAlign, 0)
	hlayoutAlign.AddLayout(vlayout, 0)

	gbNodeProperties.SetLayout(hlayoutAlign)

	hlayoutLine21Node.AddWidget(gbNodeProperties, 0, 0)

	// Incoming/Outgoing Edge Tables
	gp.EdgeIncoming = widgets.NewQTableWidget2(0, 1, nil)
	gp.EdgeOutgoing = widgets.NewQTableWidget2(0, 1, nil)

	gp.EdgeIncoming.SetVisible(false)
	gp.EdgeOutgoing.SetVisible(false)

	var strHeaders1 = []string{"Incoming Edge"}
	gp.EdgeIncoming.SetHorizontalHeaderLabels(strHeaders1)
	gp.EdgeIncoming.HorizontalHeader().SetStretchLastSection(true)

	var strHeaders2 = []string{"Outgoing Edge"}
	gp.EdgeOutgoing.SetHorizontalHeaderLabels(strHeaders2)
	gp.EdgeOutgoing.HorizontalHeader().SetStretchLastSection(true)

	hlayoutLine22Edge.AddWidget(gp.EdgeIncoming, 0, 0)
	hlayoutLine22Edge.AddWidget(gp.EdgeOutgoing, 0, 0)

	vlayoutLine2.AddLayout(hlayoutLine21Node, 0)
	vlayoutLine2.AddLayout(hlayoutLine22Edge, 0)

	var w1 = widgets.NewQWidget(nil, 0)
	var w2 = widgets.NewQWidget(nil, 0)
	var sp = widgets.NewQSplitter(nil)

	w1.SetLayout(hlayoutLine1)
	w2.SetLayout(vlayoutLine2)
	sp.AddWidget(w1)
	sp.AddWidget(w2)

	hlayoutLine3.AddWidget(sp, 0, 0)

	return hlayoutLine3

}

func (gp *GraphPanel) line3() *widgets.QVBoxLayout {

	var vlayout = widgets.NewQVBoxLayout()

	gp.TblArgs = NewArgsTable() // Args Table

	gp.ValueList = Valuelist{widgets.NewQListWidget(nil)}
	gp.ValueList.SetFont(gui.NewQFont2("verdana", 13, 1, false))
	gp.ValueList.SetAcceptDrops(true)

	gp.ValueList.ConnectDragEnterEvent(gp.ValueList.DragEnterEvent_)
	gp.ValueList.ConnectDragMoveEvent(gp.ValueList.DragMoveEvent_)
	gp.ValueList.ConnectDropEvent(gp.ValueList.DropEvent_)

	var sp = widgets.NewQSplitter(nil)
	sp.SetOrientation(2)
	sp.AddWidget(gp.TblArgs)
	sp.AddWidget(gp.ValueList)

	vlayout.AddWidget(sp, 0, 0)

	return vlayout

}
