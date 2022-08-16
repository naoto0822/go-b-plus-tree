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

func (b *BTree) Insert(key, value []byte) error {
	insertResult, err := b.insertChildNode(b.RootNode, key, value)
	if err != nil {
		return err
	}

	if insertResult.Type == InsertResultSplit {
		newRootPage, err := b.bufferPoolManager.AllocatePage(NodeTypeInternal)
		if err != nil {
			return nil
		}
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
	}

	return nil
}

func (b *BTree) insertChildNode(node Node, key []byte, value []byte) (*InsertResult, error) {
	switch node.GetNodeType() {
	case NodeTypeLeaf:
		leafNode, ok := node.(*LeafNode)
		if !ok {
			return nil, fmt.Errorf("TODO")
		}

		err := leafNode.Insert(key, value)
		if err != nil {
			return nil, err
		}

		if leafNode.Length() >= MaxOrder {
			newPage, err := b.bufferPoolManager.AllocatePage(NodeTypeLeaf)
			if err != nil {
				return nil, err
			}
			newLeafNode := NewLeafNode(newPage)

			// split records
			splitPont := leafNode.Length() / 2
			leftRecords := leafNode.Page.Records[:splitPont]
			rightRecords := leafNode.Page.Records[splitPont:]
			leafNode.Page.Records = leftRecords
			newLeafNode.Page.Records = rightRecords

			// prev <-> leafNode <-> next
			// ->
			// prev <-> leafNode <-> newLeafNode <-> next
			nextPageID := leafNode.Page.NextID
			leafNode.Page.NextID = newLeafNode.Page.ID
			newLeafNode.Page.PrevID = leafNode.Page.ID
			if nextPageID != NoSiblingPageID {
				nextPage, err := b.bufferPoolManager.FetchPage(nextPageID)
				if err != nil {
					return nil, err
				}
				nextNode := NewLeafNode(nextPage)
				newLeafNode.Page.NextID = nextNode.Page.ID
				nextNode.Page.PrevID = newLeafNode.Page.ID

				b.bufferPoolManager.Commit(nextNode.Page)
			}
			b.bufferPoolManager.Commit(leafNode.Page)
			b.bufferPoolManager.Commit(newLeafNode.Page)

			overflowKeyValue := leafNode.Page.Records[len(leafNode.Page.Records)-1]
			return &InsertResult{
				Type:           InsertResultSplit,
				OverflowKey:    overflowKeyValue.Key,
				OverflowPageID: leafNode.Page.ID,
				Left:           leafNode,
				Right:          newLeafNode,
			}, nil

		} else {
			return &InsertResult{
				Type: InsertResultTypeFit,
			}, nil
		}

	case NodeTypeInternal:
		internalNode, ok := node.(*InternalNode)
		if !ok {
			return nil, fmt.Errorf("TODO")
		}

		childPageID := internalNode.findChildPageID(key)
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

		if insertResult.Type == InsertResultSplit {
			err := internalNode.Insert(insertResult.OverflowKey, insertResult.OverflowPageID)
			if err != nil {
				return nil, err
			}

			if internalNode.Length() >= MaxOrder {
				newPage, err := b.bufferPoolManager.AllocatePage(NodeTypeInternal)
				if err != nil {
					return nil, err
				}
				newInternalNode := NewInternalNode(newPage)

				// split records
				splitPont := internalNode.Length() / 2
				leftRecords := internalNode.Page.Records[:splitPont]
				rightRecords := internalNode.Page.Records[splitPont:]
				internalNode.Page.Records = leftRecords
				newInternalNode.Page.Records = rightRecords

				b.bufferPoolManager.Commit(internalNode.Page)
				b.bufferPoolManager.Commit(newInternalNode.Page)

				overflowKeyValue := internalNode.Page.Records[len(internalNode.Page.Records)-1]
				return &InsertResult{
					Type:           InsertResultSplit,
					OverflowKey:    overflowKeyValue.Key,
					OverflowPageID: internalNode.Page.ID,
					Left:           internalNode,
					Right:          newInternalNode,
				}, nil

			} else {
				return &InsertResult{
					Type: InsertResultTypeFit,
				}, nil
			}
		} else {
			return &InsertResult{
				Type: InsertResultTypeFit,
			}, nil
		}

	default:
		return nil, fmt.Errorf("unknown node type")
	}
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
