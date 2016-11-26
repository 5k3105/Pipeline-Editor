package janus

import (
	"strconv"

	"github.com/emirpasic/gods/maps/treemap"
)

var graph *Graph

type Graph struct {
	GraphName   string
	DatasetPath string
	JobID       string
	Nodes       *treemap.Map
}

type Node struct {
	NodeID     string
	Language   string
	ScriptFile string
	ClassName  string
	Args       *treemap.Map
	Edges      *treemap.Map
	AntiEdges  *treemap.Map
	X, Y       float64
}

type Arg struct {
	Source string
	Value  string
	Values *treemap.Map
}

type Edge struct {
	String string
}

func NewJanusGraph() *Graph {
	graph = &Graph{Nodes: treemap.NewWithIntComparator()}
	return graph
}

func (jg *Graph) AddNode(i int, x, y float64) {
	jg.Nodes.Put(i, &Node{NodeID: strconv.Itoa(i),
		Language:  "python",
		ClassName: "node " + strconv.Itoa(i),
		Args:      treemap.NewWithIntComparator(),
		Edges:     treemap.NewWithIntComparator(),
		AntiEdges: treemap.NewWithIntComparator(),
		X:         x,
		Y:         y})

	in, _ := jg.Nodes.Get(i)
	n := in.(*Node)
	n.AddArg(0)
}

func (jg *Graph) AddNodeFromFile(i int, x, y float64, language, classname, scriptfile string) {
	jg.Nodes.Put(i, &Node{
		NodeID:     strconv.Itoa(i),
		Language:   language,
		ClassName:  classname,
		ScriptFile: scriptfile,
		Args:       treemap.NewWithIntComparator(),
		Edges:      treemap.NewWithIntComparator(),
		AntiEdges:  treemap.NewWithIntComparator(),
		X:          x,
		Y:          y})

	in, _ := jg.Nodes.Get(i)
	n := in.(*Node)
	n.AddArg(0)
}

func (jg *Graph) RemoveNode(i int) {
	jg.Nodes.Remove(i)

}

func (n *Node) AddEdge(i int, t string) {
	switch t {
	case "Incoming Edge":
		n.Edges.Put(i, &Edge{String: strconv.Itoa(i)}) // edge is using i as key and value
	case "Outgoing Edge":
		n.AntiEdges.Put(i, &Edge{String: strconv.Itoa(i)})
	}

}

func (n *Node) RemoveEdge(i int, t string) {
	switch t {
	case "Incoming Edge":
		n.Edges.Remove(i)
	case "Outgoing Edge":
		n.AntiEdges.Remove(i)
	}

}

func (n *Node) RemoveEdges() {

	//		n.Edges.Remove(i)

	//		n.AntiEdges.Remove(i)

}

func (n *Node) AddArg(i int) {
	n.Args.Put(i, &Arg{Source: "", Values: treemap.NewWithIntComparator()})
}

func (n *Node) RemoveArg(i int) {
	n.Args.Remove(i)
}
