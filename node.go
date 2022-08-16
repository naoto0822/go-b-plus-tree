package bplustree

import (
	"fmt"
)

type NodeType int

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeRoot
	NodeTypeInternal
	NodeTypeLeaf
)

type Node interface {
	GetNodeType() NodeType
	GetMaxKey() []byte
	GetPageID() int64
}

func NewNode(page Page) (Node, error) {
	switch page.NodeType {
	case NodeTypeInternal:
		return NewInternalNode(page), nil
	case NodeTypeLeaf:
		return NewLeafNode(page), nil
	default:
		return nil, fmt.Errorf("unknown node type")
	}
}
