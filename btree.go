package bplustree

import (
	"fmt"
)

// BTree ...
type BTree struct {
	RootNode          Node
	bufferPoolManager *BufferPoolManager
}

func NewBTree(bufferPoolManager *BufferPoolManager) *BTree {
	// TODO: fuzzy...
	root := &LeafNode{}

	return &BTree{
		RootNode:          root,
		bufferPoolManager: bufferPoolManager,
	}
}

func (b *BTree) Get(key []byte) (KeyValue, error) {
	leafNode, err := b.findLeafNode(key, b.RootNode)
	if err != nil {
		return KeyValue{}, err
	}

	keyValue, found := leafNode.Get(key)
	if !found {
		return KeyValue{}, err
	}

	return keyValue, nil
}

func (b *BTree) findLeafNode(key []byte, node Node) (*LeafNode, error) {
	// when LeafNode
	if node.GetNodeType() == NodeTypeLeaf {
		return node.(*LeafNode), nil
	}

	// when InternalNode
	internalNode, ok := node.(*InternalNode)
	if !ok {
		return nil, fmt.Errorf("not leafnode")
	}
	childPageID := internalNode.findChildPageID(key)
	page, err := b.bufferPoolManager.FetchPage(childPageID)
	if err != nil {
		return nil, err
	}
	childNode, err := NewNode(page)
	if err != nil {
		return nil, err
	}
	return b.findLeafNode(key, childNode)
}
