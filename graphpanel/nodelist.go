package graphpanel

import (
	"github.com/5k3105/FM/janus"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Nodelist struct{ *widgets.QListWidget }

func (n *Nodelist) CurrentItemChanged_(current *widgets.QListWidgetItem, previous *widgets.QListWidgetItem) {

	var di = current.Data(int(core.Qt__UserRole)).ToInt(true)
	jn, _ := janusgraph.Nodes.Get(di)
	i := jn.(*janus.Node)

	graphpanel.NodeID.SetText(i.NodeID)
	graphpanel.Language.SetText(i.Language)
	graphpanel.ScriptFile.SetText(i.ScriptFile)
	graphpanel.ClassName.SetText(i.ClassName)

	graphpanel.EdgeIncoming.ClearContents()
	graphpanel.EdgeIncoming.SetRowCount(i.AntiEdges.Size())
	graphpanel.EdgeOutgoing.ClearContents()
	graphpanel.EdgeOutgoing.SetRowCount(i.Edges.Size())

	font := gui.NewQFont2("verdana", 13, 1, false)

	for indx, x := range i.Edges.Values() {
		xi := x.(*janus.Edge)
		var wi = widgets.NewQTableWidgetItem2(xi.String, 0)
		wi.SetFont(font)
		graphpanel.EdgeOutgoing.SetItem(indx, 0, wi)
	}

	for indx, x := range i.AntiEdges.Values() {
		xi := x.(*janus.Edge)
		var wi = widgets.NewQTableWidgetItem2(xi.String, 0)
		wi.SetFont(font)
		graphpanel.EdgeIncoming.SetItem(indx, 0, wi)
	}

	// ARGS Table
	ClearAndSetModel(i) // .(*janus.Node) current

}
