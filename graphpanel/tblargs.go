package graphpanel

import (
	"github.com/5k3105/Pipeline-Editor/janus"

	"strconv"

	"github.com/emirpasic/gods/maps/treemap"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

var (
	model *core.QAbstractTableModel // QAbstractItemModel
	view  *widgets.QTableView

	textIndex   *core.QModelIndex
	currentRow  int
	currentNode *janus.Node
	currentArg  *janus.Arg
	previousRow int //?
)

type Delegate struct {
	widgets.QStyledItemDelegate //don't use *pointers or it won't work
}

type Tblargs struct{ *widgets.QTableView }

func NewArgsTable() *widgets.QTableView {
	view = widgets.NewQTableView(nil)        // setObjectName
	model = core.NewQAbstractTableModel(nil) // NewQAbstractItemModel(nil)

	delegate := InitDelegate()
	view.SetItemDelegate(delegate)
	//view.SetItemDelegateForColumn(0, delegate)
	view.SetFont(gui.NewQFont2("verdana", 13, 1, false))

	view.SetColumnWidth(1, 250)
	view.HorizontalHeader().SetStretchLastSection(true)

	view.SetSelectionBehavior(0) // SelectItems
	view.SetSelectionMode(1)     // SingleSelection

	view.ConnectKeyPressEvent(view_keypressevent)
	//view.ConnectCurrentChanged(view_currentchanged)
	view.ConnectSelectionChanged(view_selectionchanged)

	model.ConnectRowCount(model_rowcount)
	model.ConnectColumnCount(model_columncount)
	model.ConnectData(model_data)
	model.ConnectSetData(model_setdata)
	model.ConnectInsertRows(model_insertrows)
	model.ConnectFlags(model_flags)
	model.ConnectRemoveRows(model_removerows)
	model.ConnectHeaderData(model_headerdata)

	view.SetModel(model)

	return view
}

func (gp *GraphPanel) Tblargs_Clear() {
	model.BeginResetModel()
	//gp.TblArgs = NewArgsTable() //Reset()
	model.EndResetModel()
	//SetModel()

	graphpanel.ValueList.Clear()
}

func ClearAndSetModel(jn *janus.Node) { // current Node treemap - called when node in nodelist selected
	//model.Clear()
	model.BeginResetModel()
	model.EndResetModel()

	currentNode = jn
	graphpanel.ValueList.Clear()

	currentRow = 0
	a, _ := currentNode.Args.Get(currentRow)
	currentArg = a.(*janus.Arg)

	// valuelist
	font := gui.NewQFont2("verdana", 13, 1, false)
	for _, x := range currentArg.Values.Values() {

		xj := x.(string)
		n := widgets.NewQListWidgetItem2(xj, graphpanel.ValueList, 0)
		n.SetFont(font)
	}

	graphpanel.TblArgs.SelectRow(0)

}

func view_keypressevent(e *gui.QKeyEvent) {
	if e.Key() == int(core.Qt__Key_Delete) {
		if view.CurrentIndex().Row() < currentNode.Args.Size()-1 {
			// statusbar.ShowMessage("row: "+strconv.Itoa(currentRow)+"  list size: "+strconv.Itoa(listv.Size()), 0)
			model.RemoveRows(view.CurrentIndex().Row(), 1, core.NewQModelIndex())
			updateValuelist()
		}
	} else {
		view.KeyPressEventDefault(e)
	}
}

//func view_currentchanged(current *core.QModelIndex, previous *core.QModelIndex) {
//	previousRow = previous.Row()
//	currentRow = current.Row()
//}

func view_selectionchanged(selected *core.QItemSelection, deselected *core.QItemSelection) {
	updateValuelist()
}

func updateValuelist() {
	graphpanel.ValueList.Clear()
	// clear and update valuelist w/ jgraph
	// current args list row
	var r = view.CurrentIndex().Row()
	a, _ := currentNode.Args.Get(r)
	currentArg := a.(*janus.Arg)

	for _, x := range currentArg.Values.Values() {
		xj := x.(string)
		widgets.NewQListWidgetItem2(xj, graphpanel.ValueList, 0)
	}
}

func InitDelegate() *Delegate {
	item := NewDelegate(nil) //will be generated in moc.go
	item.ConnectCreateEditor(delegate_createEditor)
	item.ConnectSetEditorData(delegate_setEditorData)
	item.ConnectSetModelData(delegate_setModelData)
	item.ConnectUpdateEditorGeometry(delegate_updateEditorGeometry)
	return item
}

func delegate_createEditor(parent *widgets.QWidget, option *widgets.QStyleOptionViewItem, index *core.QModelIndex) *widgets.QWidget {
	editor := widgets.NewQLineEdit(parent)
	textIndex = index //?

	if index.Row() == currentNode.Args.Size()-1 {
		model.InsertRow(index.Row(), core.NewQModelIndex())
	}

	//editor.ConnectTextChanged(delegate_textchanged)

	return editor.QWidget_PTR()
}

// not needed??
func delegate_textchanged(text string) {
	model.SetData(textIndex, core.NewQVariant14(text), 2) // edit role = 2
}

func delegate_setEditorData(editor *widgets.QWidget, index *core.QModelIndex) {
	var value string

	//currentRow = index.Row() // needed?
	var argi, exists = currentNode.Args.Get(index.Row())

	currentArg = argi.(*janus.Arg) // needed??

	if exists {
		//var arg = argi.(*janus.Arg)
		switch index.Column() {
		case 0:
			value = currentArg.Source
		case 1:
			value = currentArg.Value
		}
	}

	lineedit := widgets.NewQLineEditFromPointer(editor.Pointer())
	lineedit.SetText(value) //.(string)
}

func delegate_setModelData(editor *widgets.QWidget, model *core.QAbstractItemModel, index *core.QModelIndex) {
	lineedit := widgets.NewQLineEditFromPointer(editor.Pointer())
	text := lineedit.Text()
	model.SetData(index, core.NewQVariant14(text), int(core.Qt__EditRole))
}

func delegate_updateEditorGeometry(editor *widgets.QWidget, option *widgets.QStyleOptionViewItem, index *core.QModelIndex) {
	editor.SetGeometry(option.Rect())
}

func model_headerdata(section int, orientation core.Qt__Orientation, role int) *core.QVariant {
	if orientation == 1 && role == 0 { // Qt__Horizontal, Qt__DisplayRole
		switch section {
		case 0:
			return core.NewQVariant14("Source")
		case 1:
			return core.NewQVariant14("Value")
		}

		return core.NewQVariant14("column" + strconv.Itoa(section+1))
	}

	if orientation == 2 && role == 0 {
		if section < currentNode.Args.Size()-1 {
			return core.NewQVariant14(strconv.Itoa(section + 1))
		} else {
			return core.NewQVariant14("*")
		}
	}
	return core.NewQVariant()
}

func model_rowcount(parent *core.QModelIndex) int {
	return currentNode.Args.Size() //+ 1
}

func model_columncount(parent *core.QModelIndex) int {
	return 2
}

func model_data(index *core.QModelIndex, role int) *core.QVariant { // dataset to model/view

	if role == 0 && index.IsValid() { // display role = 0

		var argi, exists = currentNode.Args.Get(index.Row())

		if exists {
			var arg = argi.(*janus.Arg)
			switch index.Column() {
			case 0:
				return core.NewQVariant14(arg.Source)
			case 1:
				return core.NewQVariant14(arg.Value)
			}
		}
	}
	return core.NewQVariant()

}

func model_setdata(index *core.QModelIndex, value *core.QVariant, role int) bool { // model/view to dataset
	if (role == 2 || role == 0) && index.IsValid() { // edit role = 2

		var argi, exists = currentNode.Args.Get(index.Row())

		if exists {
			var arg = argi.(*janus.Arg)
			//if !(arg.Source == "" && arg.Value == "") {
			switch index.Column() {
			case 0:
				currentNode.Args.Put(index.Row(), &janus.Arg{Source: value.ToString(), Value: arg.Value, Values: arg.Values})
			case 1:
				currentNode.Args.Put(index.Row(), &janus.Arg{Source: arg.Source, Value: value.ToString(), Values: arg.Values})
			}
			return true
		}
	}
	return true
}

func model_insertrows(row int, count int, parent *core.QModelIndex) bool {
	model.BeginInsertRows(core.NewQModelIndex(), row, row)

	currentNode.Args.Put(row+1, &janus.Arg{Source: "", Values: treemap.NewWithIntComparator()})

	model.EndInsertRows()
	view.SelectRow(row)
	return true
}

func model_removerows(row int, count int, parent *core.QModelIndex) bool {
	model.BeginRemoveRows(core.NewQModelIndex(), row, row)

	currentNode.Args.Remove(row)

	// reset indexes
	for i, x := range currentNode.Args.Keys() { // i always starts from zero
		indx := x.(int) // indx does not

		if indx != i {
			var ar, _ = currentNode.Args.Get(indx)
			var arg = ar.(*janus.Arg)
			currentNode.Args.Put(i, &janus.Arg{Source: arg.Source, Value: arg.Value, Values: arg.Values})
			currentNode.Args.Remove(indx)
		}

	}

	model.EndRemoveRows()
	return true

}

func model_flags(index *core.QModelIndex) core.Qt__ItemFlag {
	return 35 // 1 || 2 || 32 // ItemIsSelectable || ItemIsEditable || ItemIsEnabled
}

//func (ta *Tblargs) ItemChanged_(current *widgets.QTableWidgetItem) {
//	var di = graphpanel.NodeList.CurrentItem().Data(int(core.Qt__UserRole)).ToInt(true)
//	i, _ := janusgraph.Nodes.Get(di)
//	jn := i.(*janus.Node)

//	wit := current.Text()
//	currentRow := current.Row()
//	currentColumn := current.Column()

//	for indx, x := range jn.Args.Values() {
//		xi := x.(*janus.Arg)

//		if currentRow == indx {
//			switch currentColumn {
//			case 1: // Source
//				xi.Source = wit

//			case 2: // Value
//				xi.Value = wit

//			}
//		}

//	}

//}

//func (ta *Tblargs) CurrentCellChanged(currentRow int, currentColumn int, previousRow int, previousColumn int) {
//	var di = graphpanel.NodeList.CurrentItem().Data(int(core.Qt__UserRole)).ToInt(true)
//	i, _ := janusgraph.Nodes.Get(di)
//	jn := i.(*janus.Node)

//	wit := ta.Item(currentRow, currentColumn).Text()

//	for indx, x := range jn.Args.Values() {
//		xi := x.(*janus.Arg)

//		if currentRow == indx {
//			switch currentColumn {
//			case 1: // Source
//				xi.Source = wit

//			case 2: // Value
//				xi.Value = wit

//			}
//		}

//	}

//}

//		var wi = widgets.NewQTableWidgetItem2(xi.Source, 0)
//		wi.SetFont(font)
//		graphpanel.TblArgs.SetItem(indx, 1, wi)

//		for _, x := range xi.Value.Values() {
//			xi := x.(string)

//			if xi != "" {
//				n := widgets.NewQListWidgetItem2(xi, graphpanel.ValueList, 0)
//				n.SetFont(font)
//			}

//		}

//https://github.com/RustamSafiulin/mnemo_designer/blob/478e34be3564d3c51d8fae7aac9b3c820c6bc12c/propertyeditor/propertytableview.cpp
//void PropertyTableView::init() {

//    resizeColumnsToContents();
//    setColumnWidth(0,150);
//    setAlternatingRowColors(true);
//    setSelectionMode(QTableView::SingleSelection);
//    setSelectionBehavior(QTableView::SelectRows);
//    setEditTriggers(QAbstractItemView::CurrentChanged | QAbstractItemView::SelectedClicked);
//    verticalHeader()->hide();
//    verticalHeader()->setDefaultSectionSize(25);
//    horizontalHeader()->setStretchLastSection(true);

//    delegate = new PropertyDelegate(this);
//    setItemDelegate(delegate);
//}
