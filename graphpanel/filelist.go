package graphpanel

import (
	"strings"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

type Filelist struct{ *widgets.QListWidget }

func (f *Filelist) KeyPressEvent(e *gui.QKeyEvent) {

	if e.Key() == int(core.Qt__Key_Delete) {
		f.TakeItem(f.CurrentRow())
	} else {
		view.KeyPressEventDefault(e)
	}

}

func (f *Filelist) CurrentItemChanged(current *widgets.QListWidgetItem, previous *widgets.QListWidgetItem) {
	txt := current.Text()
	LoadXML(txt)
	graphpanel.FileName.SetText(txt)
}

func (f *Filelist) DragEnterEvent_(e *gui.QDragEnterEvent) {
	e.AcceptProposedAction()
}

func (f *Filelist) DragMoveEvent_(e *gui.QDragMoveEvent) {
	e.AcceptProposedAction()
}

func (f *Filelist) DropEvent_(e *gui.QDropEvent) {
	e.SetDropAction(core.Qt__CopyAction)
	e.AcceptProposedAction()
	e.SetAccepted(true)

	font := gui.NewQFont2("verdana", 13, 1, false)

	for _, j := range strings.Split(e.MimeData().Text(), "\n") { // Urls()
		if j[len(j)-4:] == ".xml" {
			n := widgets.NewQListWidgetItem2(j[8:], f, 0)
			n.SetFont(font)

		}
	}

}
