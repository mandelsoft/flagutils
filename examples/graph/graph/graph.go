package graph

type Node struct {
	name     string
	value    string
	children []*Node
}

func NewNode(name string, value string) *Node {
	return &Node{
		name:  name,
		value: value,
	}
}

func (n *Node) Name() string {
	return n.name
}

func (n *Node) Value() string {
	return n.value
}

func (n *Node) HasChildren() bool {
	return len(n.children) > 0
}

func (n *Node) AddChild(child *Node) {
	n.children = append(n.children, child)
}

func (n *Node) Children(yield func(*Node) bool) {
	for _, child := range n.children {
		if !yield(child) {
			return
		}
	}
}

////////////////////////////////////////////////////////////////////////////////

func NewGraph() *Graph {
	return &Graph{}
}

type Graph struct {
	roots []*Node
}

func (g *Graph) AddRoot(node ...*Node) *Graph {
	g.roots = append(g.roots, node...)
	return g
}

func (g *Graph) Roots(yield func(*Node) bool) {
	for _, root := range g.roots {
		if !yield(root) {
			return
		}
	}
}

func (g *Graph) GetRoot(name string) *Node {
	for _, root := range g.roots {
		if root.Name() == name {
			return root
		}
	}
	return nil
}
