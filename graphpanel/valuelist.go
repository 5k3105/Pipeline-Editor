package graphpanel

import (
	"github.com/5k3105/FM/janus"
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Valuelist struct{ *widgets.QListWidget }

func (v *Valuelist) DragEnterEvent_(e *gui.QDragEnterEvent) {
	e.AcceptProposedAction()
}

func (v *Valuelist) DragMoveEvent_(e *gui.QDragMoveEvent) {
	e.AcceptProposedAction()
}

func (v *Valuelist) DropEvent_(e *gui.QDropEvent) {
	e.SetDropAction(core.Qt__CopyAction)
	e.AcceptProposedAction()
	e.SetAccepted(true)

	font := gui.NewQFont2("verdana", 13, 1, false)

	// current node
	var di = graphpanel.NodeList.CurrentItem().Data(int(core.Qt__UserRole)).ToInt(true)
	ji, _ := janusgraph.Nodes.Get(di)
	jn := ji.(*janus.Node)

	// current args list row
	var r = view.CurrentIndex().Row()
	a, _ := jn.Args.Get(r)

	ja := a.(*janus.Arg)
	ja.Values.Clear()

	// write list items and j.Arg.Values
	for indx, j := range strings.Split(e.MimeData().Text(), "\n") { // Urls()
		if j != "" {
			n := widgets.NewQListWidgetItem2(j[8:], v, 0)
			n.SetFont(font)

			ja.Values.Put(indx, j[8:])

		}
	}

}
