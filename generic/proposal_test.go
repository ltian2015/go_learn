package generic

// 本文件中的代码来自Golang的泛型提案中的代码。
// https://go.googlesource.com/proposal/+/HEAD/design/43651-type-parameters.md#Type-parameters
// NodeConstraint is the type constraint for graph nodes:
// they must have an Edges method that returns the Edge's
// that connect to this Node.
type NodeConstraint[Edge any] interface {
	Edges() []Edge
}

// EdgeConstraint is the type constraint for graph edges:
// they must have a Nodes method that returns the two Nodes
// that this edge connects.
type EdgeConstraint[Node any] interface {
	Nodes() (from, to Node)
}

// Graph is a graph composed of nodes and edges.
type Graph[Node NodeConstraint[Edge], Edge EdgeConstraint[Node]] struct{}

// New returns a new graph given a list of nodes.
func New[Node NodeConstraint[Edge], Edge EdgeConstraint[Node]](nodes []Node) *Graph[Node, Edge] {
	return nil
}

// ShortestPath returns the shortest path between two nodes,
// as a list of edges.
func (g *Graph[Node, Edge]) ShortestPath(from, to Node) []Edge {
	return nil
}

// Vertex is a node in a graph.
type Vertex struct{}

// Edges returns the edges connected to v.
func (v *Vertex) Edges() []*FromTo {
	return nil
}

// FromTo is an edge in a graph.
type FromTo struct{}

// Nodes returns the nodes that ft connects.
func (ft *FromTo) Nodes() (*Vertex, *Vertex) {
	return nil, nil
}

var g = New([]*Vertex{})

type NodeInterface interface {
	Edges() []EdgeInterface
}
type EdgeInterface interface {
	Nodes() (NodeInterface, NodeInterface)
}

var node NodeConstraint[*FromTo] = &Vertex{}
