package genxml

import (
	"github.com/5k3105/FM/gfxcanvas"
	"github.com/5k3105/FM/gfxinterface"
	"github.com/5k3105/FM/graphpanel"
	"github.com/5k3105/FM/janus"
	"strconv"

	"github.com/beevik/etree"
	"github.com/emirpasic/gods/maps/treemap"
)

const Scalefactor = 10

var w, h float64 = 30 * Scalefactor, 10 * Scalefactor

func LoadXML(filename string, canvas *gfxcanvas.Canvas, gp *graphpanel.GraphPanel) { //}*janus.Graph {
	canvas.Statusbar.ShowMessage(" start loadxml ", 0)

	gp.SetJanusGraph(canvas.JanusGraph) // update pointer to new janusgraph

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(filename); err != nil {
		panic(err)
	}

	graph := doc.SelectElement("graph")

	// update janusgraph
	canvas.JanusGraph.GraphName = graph.SelectElement("graphName").Text()
	canvas.JanusGraph.DatasetPath = graph.SelectElement("datasetPath").Text()
	canvas.JanusGraph.JobID = graph.SelectElement("jobID").Text()

	// update graphpanel
	gp.GraphName.SetText(graph.SelectElement("graphName").Text())
	gp.DatasetPath.SetText(graph.SelectElement("datasetPath").Text())
	gp.JobID.SetText(graph.SelectElement("jobID").Text())

	nodes := graph.SelectElement("nodes")

	// create all nodes
	for _, node := range nodes.SelectElements("node") {

		x, _ := strconv.ParseFloat(node.SelectElement("x").Text(), 64)
		y, _ := strconv.ParseFloat(node.SelectElement("y").Text(), 64)
		i, _ := strconv.Atoi(node.SelectElement("nodeID").Text())

		canvas.AddRectangleFromFile(i, x, y, w, h,
			node.SelectElement("language").Text(),
			node.SelectElement("className").Text(),
			node.SelectElement("scriptFile").Text())
	}

	// create all edges
	for _, node := range nodes.SelectElements("node") {

		sid, _ := strconv.Atoi(node.SelectElement("nodeID").Text())
		s, _ := canvas.Figures.Get(sid)
		source := s.(gfxinterface.Figure)

		edges := node.SelectElement("edges")

		for _, edge := range edges.SelectElements("edge") {
			//if edge != struct{} {
			tid, _ := strconv.Atoi(edge.Text())
			t, _ := canvas.Figures.Get(tid)
			target := t.(gfxinterface.Figure)
			canvas.AddLineFromFile(source, target)

		}
	}

	for _, node := range nodes.SelectElements("node") {

		sid, _ := strconv.Atoi(node.SelectElement("nodeID").Text())
		ji, _ := canvas.JanusGraph.Nodes.Get(sid)
		jn := ji.(*janus.Node)

		args := node.SelectElement("args")

		for i, arg := range args.ChildElements() {

			if arg.SelectElement("source").Text() == "" {
				return
			}

			if len(arg.SelectElements("value")) > 1 {
				values := treemap.NewWithIntComparator()
				for n, ar := range arg.SelectElements("value") {
					values.Put(n, ar.Text())
				}
				jn.Args.Put(i, &janus.Arg{Source: arg.SelectElement("source").Text(),
					Values: values})
			} else {
				jn.Args.Put(i, &janus.Arg{Source: arg.SelectElement("source").Text(),
					Value:  arg.SelectElement("value").Text(),
					Values: treemap.NewWithIntComparator()})
			}
		}

		// extra element for args table edit-line
		jn.Args.Put(len(args.ChildElements()), &janus.Arg{Source: "",
			Value:  "",
			Values: treemap.NewWithIntComparator()})

	}

	canvas.Statusbar.ShowMessage(" end loadxml ", 0)
}

func GenXML(gp *graphpanel.GraphPanel, jg *janus.Graph, filename string) { // gp to get graphname etc

	doc := etree.NewDocument()
	doc.CreateProcInst("xml", `version="1.0" encoding="UTF-8"`)
	//doc.CreateProcInst("xml-stylesheet", `type="text/xsl" href="style.xsl"`)

	graph := doc.CreateElement("graph")

	jg.GraphName = gp.GraphName.Text()
	jg.DatasetPath = gp.DatasetPath.Text()
	jg.JobID = gp.JobID.Text()

	graph.CreateElement("graphName").SetText(jg.GraphName)
	graph.CreateElement("datasetPath").SetText(jg.DatasetPath)
	graph.CreateElement("jobID").SetText(jg.JobID)

	nodes := graph.CreateElement("nodes")

	for _, x := range jg.Nodes.Values() {
		nv := x.(*janus.Node)
		node := nodes.CreateElement("node")

		node.CreateElement("nodeID").SetText(nv.NodeID)
		node.CreateElement("language").SetText(nv.Language)
		node.CreateElement("scriptFile").SetText(nv.ScriptFile)
		node.CreateElement("className").SetText(nv.ClassName)

		node.CreateElement("x").SetText(FloatToString(nv.X))
		node.CreateElement("y").SetText(FloatToString(nv.Y))

		args := node.CreateElement("args")

		for i, x := range nv.Args.Values() {
			av := x.(*janus.Arg)
			if av.Source != "" {
				arg := args.CreateElement("arg" + strconv.Itoa(i))

				arg.CreateElement("source").SetText(av.Source)
				//if av.Value != "" {
				arg.CreateElement("value").SetText(av.Value)
				//}

				if !av.Values.Empty() {
					for _, x := range av.Values.Values() {
						xs := x.(string)
						arg.CreateElement("value").SetText(xs)
					}
				}
			}
		}

		edges := node.CreateElement("edges")

		for _, x := range nv.Edges.Values() {
			xv := x.(*janus.Edge)
			edges.CreateElement("edge").SetText(xv.String)
		}

		antiedges := node.CreateElement("antiEdges")

		for _, x := range nv.AntiEdges.Values() {
			xv := x.(*janus.Edge)
			antiedges.CreateElement("antiEdge").SetText(xv.String)
		}

	}

	doc.Indent(2)
	doc.WriteToFile(filename)

}

func FloatToString(input_num float64) string {
	return strconv.FormatFloat(input_num, 'f', 1, 64)
}
