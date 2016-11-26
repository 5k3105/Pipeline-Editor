package point

import (
	"github.com/5k3105/Pipeline-Editor/gfxinterface"
	"github.com/5k3105/Pipeline-Editor/graph"
)

type Point struct {
	X, Y float64
}

func (p *Point) GetX() float64        { return p.X }
func (p *Point) SetX(x float64)       { p.X = x }
func (p *Point) GetY() float64        { return p.Y }
func (p *Point) SetY(y float64)       { p.Y = y }
func (p *Point) GetW() float64        { return 0.0 }
func (p *Point) GetNode() *graph.Node { return nil }

func (p *Point) AddEdgeIncoming(edge gfxinterface.Link) {}

func (p *Point) AddEdgeOutgoing(edge gfxinterface.Link) {}
