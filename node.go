package bplustree

type NodeType int

const (
	NodeTypeRoot NodeType = iota
	NodeTypeInternal
	NodeTypeLeaf
)
