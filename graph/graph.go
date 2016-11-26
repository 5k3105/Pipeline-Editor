package graph

import (
	"encoding/json"
	"os"
)

var autonumber int

type Graph struct {
	Id    int
	Nodes []*Node `json:",omitempty"`
}

type Node struct {
	Id            int
	IncomingEdges []*Edge
	OutgoingEdges []*Edge
}

type Edge struct {
	Id     int
	Source *Node `json:",omitempty"`
	Target *Node `json:",omitempty"`
}

func NewGraph() *Graph {
	var g Graph
	g.Id = AutoNumber()
	return &g
}

// Encode to file
func Save(path string, g *Graph) error {
	file, err := os.Create(path)
	if err == nil {
		encoder := json.NewEncoder(file)
		encoder.Encode(g)
	}
	file.Close()
	return err
}

// Decode file
func Load(path string, g *Graph) (error, *Graph) {
	file, err := os.Open(path)
	if err == nil {
		decoder := json.NewDecoder(file)
		err = decoder.Decode(g)
	}
	file.Close()
	return err, g
}

func (g *Graph) AddNode() *Node {
	var n Node

	n.Id = AutoNumber()
	g.Nodes = append(g.Nodes, &n)

	return &n
}

func (g *Graph) AddNodeFromFile(i int) *Node {
	var n Node

	n.Id = i       //AutoNumber()
	autonumber = i // hack
	g.Nodes = append(g.Nodes, &n)

	return &n
}

func (g *Graph) AddEdge(source, target *Node) *Edge {
	var e Edge

	e.Id = AutoNumber()
	e.Source = source
	e.Target = target

	s, er := g.NodeFind(source.Id)
	Check(er)

	t, er := g.NodeFind(source.Id)
	Check(er)

	g.Nodes[s].OutgoingEdges = append(g.Nodes[s].OutgoingEdges, &e)
	g.Nodes[t].IncomingEdges = append(g.Nodes[t].IncomingEdges, &e)

	return &e
}

func (g *Graph) NodeFind(Id int) (i int, err error) {
	var s = g.Nodes

	for i, v := range s {
		if v.Id == Id {
			return i, err
		}
	}
	return 0, err
}

func AutoNumber() int {
	autonumber += 1
	return autonumber
}

func Check(e error) {
	if e != nil {
		//_, file, line, _ := runtime.Caller(1)
		//fmt.Println(line, "\t", file, "\n", e)
		//os.Exit(1)
	}
}

//func AutoNumber() Id {
//	var i Id
//	autonumber += 1
//	i = autonumber
//	return i
//}

//func (g *Graph) RemoveNode(n Node) []Node {
//	var s = g.Nodes

//	for i, v := range s {
//		if v == n {
//			g.Nodes = append(s[0:i], s[i+1:]...)
//			return g.Nodes
//		}
//	}
//	return s

//}

//func (g *Graph) RemoveEdge() {

//}

//------------------------------------------------------------------------------

//type Node struct {
//	Id
//	//X, Y float64 `json:",omitempty"`
//	//Incoming []Edge  `json:",omitempty"`
//	//Outgoing []Edge  `json:",omitempty"`
//}

//func (g *Graph) RemoveNode(n *Node) {

//}

//func (g *Graph) RemoveEdge() {

//}
