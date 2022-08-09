package bplustree

type NodeType int

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeRoot
	NodeTypeInternal
	NodeTypeLeaf
)

type Node interface {
	GetNodeType() NodeType
	IsRootNode() bool
	IsInternalNode() bool
	IsLeafNode() bool
}
