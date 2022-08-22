package bplustree

import (
	"fmt"
)

type NodePrinter struct{}

func (n *NodePrinter) Print(node Node, bufferPoolManager *BufferPoolManager) error {
	return n.print(node, bufferPoolManager)
}

func (n *NodePrinter) print(node Node, bufferPoolManager *BufferPoolManager) error {
	g := newGraph()

	if node.GetNodeType() == NodeTypeLeaf {
		g.addEdge(node.String(), "", "")
		fmt.Println(g)
		return nil
	}

	q := queue{}
	q.push(node)

	for !q.empty() {
		parent, err := q.pop()
		if err != nil {
			return err
		}

		switch parent.GetNodeType() {
		case NodeTypeInternal:
			for _, r := range parent.GetRecords() {
				pageID := decodePageID(r.Value)
				page, err := bufferPoolManager.FetchPage(pageID)
				if err != nil {
					return err
				}
				child, err := NewNode(page)
				if err != nil {
					return err
				}

				g.addEdge(parent.String(), child.String(), "")
				q.push(child)
			}

		case NodeTypeLeaf:
			continue

		default:
			return fmt.Errorf("Unknown NodeType: %v", parent.GetNodeType())
		}
	}

	fmt.Println(g)
	return nil
}

type edge struct {
	node  string
	label string
}

type graph struct {
	nodes map[string][]edge
}

func newGraph() *graph {
	return &graph{nodes: make(map[string][]edge)}
}

func (g *graph) addEdge(from, to, label string) {
	g.nodes[from] = append(g.nodes[from], edge{node: to, label: label})
}

func (g *graph) getEdges(node string) []edge {
	return g.nodes[node]
}

func (e *edge) String() string {
	return fmt.Sprintf("%v", e.node)
}

// TODO: when single leaf
func (g *graph) String() string {
	out := `digraph "" {
		graph [fontname = "Ricty Bold"];
    	node [fontname = "Ricty Bold", label="\N"];
    	edge [fontname = "Ricty Bold"];`
	for k := range g.nodes {
		for _, v := range g.getEdges(k) {
			out += fmt.Sprintf("\t\"%s\" -> \"%s\"\t[ label = \"%s\" ];\n", k, v.node, v.label)
		}
	}
	out += "}"
	return out
}

type queue struct {
	Data []Node
}

func (q *queue) push(node Node) {
	q.Data = append(q.Data, node)
}

func (q *queue) pop() (Node, error) {
	if len(q.Data) == 0 {
		return nil, fmt.Errorf("Data is emptyr")
	}
	pop := q.Data[0]
	if len(q.Data) > 1 {
		q.Data = q.Data[1:]
	} else {
		q.Data = make([]Node, 0, 0)
	}
	return pop, nil
}

func (q *queue) empty() bool {
	return len(q.Data) == 0
}
