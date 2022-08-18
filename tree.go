package bplustree

import (
	"fmt"
)

const (
	MaxOrder = 5

	PageSize = 512

	NoSiblingPageID = -1
)

// Tree ...
type Tree struct {
	RootNode          Node
	bufferPoolManager *BufferPoolManager
}

// NewTree ...
func NewTree(bufferPoolManager *BufferPoolManager) *Tree {
	rootPage := bufferPoolManager.AllocatePage(NodeTypeLeaf)
	root := NewLeafNode(rootPage)

	return &Tree{
		RootNode:          root,
		bufferPoolManager: bufferPoolManager,
	}
}

// Get ...
func (b *Tree) Get(key []byte) (KeyValue, error) {
	leafNode, err := b.findLeafNode(key, b.RootNode)
	if err != nil {
		return KeyValue{}, err
	}

	keyValue, found := leafNode.Get(key)
	if !found {
		return KeyValue{}, fmt.Errorf("TODO")
	}

	return keyValue, nil
}

// Insert ...
func (b *Tree) Insert(key, value []byte) error {
	insertResult, err := b.insertChildNode(b.RootNode, key, value)
	if err != nil {
		return err
	}

	if insertResult.Type == InsertResultSplit {
		newRootPage := b.bufferPoolManager.AllocatePage(NodeTypeInternal)
		newRootNode := NewInternalNode(newRootPage)

		err = newRootNode.Insert(insertResult.Left.GetMaxKey(), insertResult.Left.GetPageID())
		if err != nil {
			return err
		}
		err = newRootNode.Insert(insertResult.Right.GetMaxKey(), insertResult.Right.GetPageID())
		if err != nil {
			return err
		}

		b.RootNode = newRootNode
		b.bufferPoolManager.Flush(newRootNode.Page)
	}

	return nil
}

func (b *Tree) insertChildNode(node Node, key []byte, value []byte) (*InsertResult, error) {
	switch node.GetNodeType() {
	case NodeTypeLeaf:
		leafNode, ok := node.(*LeafNode)
		if !ok {
			return nil, fmt.Errorf("TODO")
		}

		isOverMaxKey := leafNode.IsOverMaxKey(key)

		err := leafNode.Insert(key, value)
		if err != nil {
			return nil, err
		}

		if leafNode.Length() >= MaxOrder {
			newPage := b.bufferPoolManager.AllocatePage(NodeTypeLeaf)
			newLeafNode := NewLeafNode(newPage)

			// split records
			// leafNode
			// ->
			// (newLeafNode, leafNode)
			splitPont := leafNode.Length() / 2

			leftRecords := make([]KeyValue, splitPont)
			copy(leftRecords, leafNode.Page.Records[:splitPont])
			rightRecords := make([]KeyValue, leafNode.Length()-splitPont)
			copy(rightRecords, leafNode.Page.Records[splitPont:])
			newLeafNode.Page.Records = leftRecords
			leafNode.Page.Records = rightRecords

			// prev <-> leafNode <-> next
			// ->
			// prev <-> newLeafNode <-> leafNode <-> next
			prevPageID := leafNode.Page.PrevID
			newLeafNode.Page.NextID = leafNode.Page.ID
			leafNode.Page.PrevID = newLeafNode.Page.ID
			if prevPageID != NoSiblingPageID {
				prevPage, err := b.bufferPoolManager.FetchPage(prevPageID)
				if err != nil {
					return nil, err
				}
				prevNode := NewLeafNode(prevPage)
				prevNode.Page.NextID = newLeafNode.Page.ID
				newLeafNode.Page.PrevID = prevNode.Page.ID
				b.bufferPoolManager.Flush(prevNode.Page)
			}

			b.bufferPoolManager.Flush(newLeafNode.Page)
			b.bufferPoolManager.Flush(leafNode.Page)

			return &InsertResult{
				Type:         InsertResultSplit,
				Left:         newLeafNode,
				Right:        leafNode,
				IsOverMaxKey: isOverMaxKey,
			}, nil

		} else {
			b.bufferPoolManager.Flush(leafNode.Page)

			return &InsertResult{
				Type:         InsertResultTypeFit,
				IsOverMaxKey: isOverMaxKey,
			}, nil
		}

	case NodeTypeInternal:
		internalNode, ok := node.(*InternalNode)
		if !ok {
			return nil, fmt.Errorf("TODO")
		}

		childPageID, found := internalNode.FindChildPageID(key)
		if !found {
			return nil, fmt.Errorf("TODO")
		}

		childPage, err := b.bufferPoolManager.FetchPage(childPageID)
		if err != nil {
			return nil, err
		}
		childNode, err := NewNode(childPage)
		if err != nil {
			return nil, err
		}

		insertResult, err := b.insertChildNode(childNode, key, value)
		if err != nil {
			return nil, err
		}

		isOverMaxKey := false
		if insertResult.IsOverMaxKey {
			isOverMaxKey = true
			internalNode.UpdateMaxKey(childNode.GetMaxKey(), childNode.GetPageID())
		}

		if insertResult.Type == InsertResultSplit {
			err := internalNode.Insert(insertResult.GetOverflowKey(), insertResult.GetOverflowPageID())
			if err != nil {
				return nil, err
			}

			if internalNode.Length() >= MaxOrder {
				newPage := b.bufferPoolManager.AllocatePage(NodeTypeInternal)
				newInternalNode := NewInternalNode(newPage)

				// split records
				// internalNode
				// ->
				// (newInternalNode, internalNode)
				splitPont := internalNode.Length() / 2

				leftRecords := make([]KeyValue, splitPont)
				copy(leftRecords, internalNode.Page.Records[:splitPont])
				rightRecords := make([]KeyValue, internalNode.Length()-splitPont)
				copy(rightRecords, internalNode.Page.Records[splitPont:])
				newInternalNode.Page.Records = leftRecords
				internalNode.Page.Records = rightRecords

				b.bufferPoolManager.Flush(newInternalNode.Page)
				b.bufferPoolManager.Flush(internalNode.Page)

				return &InsertResult{
					Type:         InsertResultSplit,
					Left:         newInternalNode,
					Right:        internalNode,
					IsOverMaxKey: isOverMaxKey,
				}, nil

			} else {
				b.bufferPoolManager.Flush(internalNode.Page)

				return &InsertResult{
					Type:         InsertResultTypeFit,
					IsOverMaxKey: isOverMaxKey,
				}, nil
			}
		} else {
			b.bufferPoolManager.Flush(internalNode.Page)

			return &InsertResult{
				Type:         InsertResultTypeFit,
				IsOverMaxKey: isOverMaxKey,
			}, nil
		}

	default:
		return nil, fmt.Errorf("unknown node type")
	}
}

func (b *Tree) findLeafNode(key []byte, node Node) (*LeafNode, error) {
	// when LeafNode
	if node.GetNodeType() == NodeTypeLeaf {
		return node.(*LeafNode), nil
	}

	// when InternalNode
	internalNode, ok := node.(*InternalNode)
	if !ok {
		return nil, fmt.Errorf("not leafnode")
	}
	childPageID, found := internalNode.FindChildPageID(key)
	if !found {
		return nil, fmt.Errorf("not found LeafNode")
	}
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

// Print ...
func (b *Tree) Print() error {
	nodePrinter := &NodePrinter{}
	err := nodePrinter.Print(b.RootNode, b.bufferPoolManager)
	if err != nil {
		return err
	}
	return nil
}
