package bplustree

import (
	"bytes"
	"fmt"

	graphviz "github.com/goccy/go-graphviz"
)

type NodePrinter struct{}

func (n *NodePrinter) Print(node Node, bufferPoolManager *BufferPoolManager) error {
	return goccyPrint(node, bufferPoolManager)
}

func originalPrint(node Node, bufferPoolManager *BufferPoolManager) error {
	g := newGraph()

	q := Queue{}
	q.Push(node)

	for !q.Empty() {
		pop, err := q.Pop()
		if err != nil {
			return err
		}

		switch pop.GetNodeType() {
		case NodeTypeInternal:
			for _, r := range pop.GetRecords() {
				pageID := decodePageID(r.Value)
				page, err := bufferPoolManager.FetchPage(pageID)
				if err != nil {
					return err
				}
				node, err := NewNode(page)
				if err != nil {
					return err
				}

				g.addEdge(pop.String(), node.String(), "")

				q.Push(node)
			}
		case NodeTypeLeaf:
			continue
		default:
			return fmt.Errorf("TODO")
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

func goccyPrint(node Node, bufferPoolManager *BufferPoolManager) error {
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		return err
	}
	defer func() {
		graph.Close()
		g.Close()
	}()

	q := Queue{}
	q.Push(node)

	dedup := make(map[int64]struct{})

	for !q.Empty() {
		pop, err := q.Pop()
		if err != nil {
			return err
		}

		_, ok := dedup[pop.GetPageID()]
		if ok {
			continue
		}

		parent, _ := graph.CreateNode(pop.String())
		dedup[pop.GetPageID()] = struct{}{}

		switch pop.GetNodeType() {
		case NodeTypeInternal:
			for _, r := range pop.GetRecords() {
				pageID := decodePageID(r.Value)
				page, err := bufferPoolManager.FetchPage(pageID)
				if err != nil {
					return err
				}
				node, err := NewNode(page)
				if err != nil {
					return err
				}
				child, _ := graph.CreateNode(node.String())
				graph.CreateEdge("hoge", parent, child)

				q.Push(node)
			}
		case NodeTypeLeaf:
			continue
		default:
			return fmt.Errorf("TODO")
		}
	}

	if err := g.RenderFilename(graph, graphviz.JPG, "./graph.png"); err != nil {
		return err
	}

	var buf bytes.Buffer
	if err := g.Render(graph, "dot", &buf); err != nil {
		return err
	}
	fmt.Println(buf.String())

	return nil
}

type Queue struct {
	Data []Node
}

func (q *Queue) Push(node Node) {
	q.Data = append(q.Data, node)
}

func (q *Queue) Pop() (Node, error) {
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

func (q *Queue) Empty() bool {
	return len(q.Data) == 0
}
