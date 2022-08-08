package bplustree

type NodeType int

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeRoot
	NodeTypeInternal
	NodeTypeLeaf
)
