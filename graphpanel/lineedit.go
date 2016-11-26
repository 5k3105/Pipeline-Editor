package graphpanel

import (
	"github.com/5k3105/Pipeline-Editor/janus"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

type Lineedit struct{ *widgets.QLineEdit }

func (l *Lineedit) EditingFinished() {

	var di = graphpanel.NodeList.CurrentItem().Data(int(core.Qt__UserRole)).ToInt(true)

	i, _ := janusgraph.Nodes.Get(di)
	jn := i.(*janus.Node)

	switch l.ObjectName() {

	case "GraphName":
		janusgraph.GraphName = l.Text()

	case "JobID":
		janusgraph.JobID = l.Text()

	case "DatasetPath":
		janusgraph.DatasetPath = l.Text()

	case "NodeID":
		jn.NodeID = l.Text()

	case "Language":
		jn.Language = l.Text()

	case "ScriptFile":
		jn.ScriptFile = l.Text()

	case "ClassName":
		jn.ClassName = l.Text()
		// update canvas and list names
		graphpanel.NodeList.CurrentItem().SetText(l.Text())
		//		f, _ := canvas.Figures.Get(di)
		//		r := f.(*rectangle.Rectangle)
		//		r.Name = l.Text()
	}

}
