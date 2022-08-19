package bplustree

import (
	"bytes"
	"fmt"
)

type NodeType int

const (
	NodeTypeUnknown NodeType = iota
	NodeTypeRoot
	NodeTypeInternal
	NodeTypeLeaf
)

// Node ...
type Node interface {
	GetNodeType() NodeType
	GetMaxKey() []byte
	GetPageID() int64
	GetRecords() []KeyValue
	String() string
}

// NewNode ...
func NewNode(page *Page) (Node, error) {
	switch page.NodeType {
	case NodeTypeInternal:
		return NewInternalNode(page), nil
	case NodeTypeLeaf:
		return NewLeafNode(page), nil
	default:
		return nil, fmt.Errorf("unknown node type")
	}
}

type FindResultType int

const (
	FindResultTypeUnknown FindResultType = iota
	FindResultTypeNoRecord
	FindResultTypeMatch
	FindResultTypeFirstGraterThanMatch
	FindResultTypeOver
)

// FindResult ...
type FindResult struct {
	Type     FindResultType
	Index    int64
	KeyValue KeyValue
}

// BaseNode ...
type BaseNode struct{}

func (b BaseNode) find(childrens []KeyValue, key []byte) *FindResult {
	if len(childrens) == 0 {
		return &FindResult{
			Type: FindResultTypeNoRecord,
		}
	}

	for idx, record := range childrens {
		compared := bytes.Compare(record.Key, key)
		switch compared {
		// Match
		// record.Key = key
		case 0:
			return &FindResult{
				Type:     FindResultTypeMatch,
				Index:    int64(idx),
				KeyValue: record,
			}

		// FirstGreaterThanMatch
		// record.Key > key
		case 1:
			return &FindResult{
				Type:  FindResultTypeFirstGraterThanMatch,
				Index: int64(idx),
				// KeyValue: , TODO: record is maybe ok
			}
		}
	}

	// Over
	return &FindResult{
		Type: FindResultTypeOver,
	}
}

func (b BaseNode) rangeScan(childrens []KeyValue, startKey, endKey []byte) ([]KeyValue, bool) {
	if len(childrens) == 0 {
		return nil, false
	}

	hits := make([]KeyValue, 0, len(childrens))
	isLastHit := false

	for idx, record := range childrens {
		startCmp := bytes.Compare(record.Key, startKey)
		endCmp := bytes.Compare(record.Key, endKey)

		switch {
		// startKey == record
		case startCmp == 0:
			hits = append(hits, record)
			if (len(childrens) - 1) == idx {
				isLastHit = true
			}

		// startKey < record < endKey
		case startCmp == 1 && endCmp == -1:
			hits = append(hits, record)
			if (len(childrens) - 1) == idx {
				isLastHit = true
			}

		// endKey == record
		case endCmp == 0:
			hits = append(hits, record)
			if (len(childrens) - 1) == idx {
				isLastHit = true
			}

		default:
			continue
		}
	}

	return hits, isLastHit
}

type InsertResultType int

const (
	InsertResultTypeUnknown InsertResultType = iota
	InsertResultTypeFit
	InsertResultSplit
)

// InsertResult ...
type InsertResult struct {
	Type         InsertResultType
	Left         Node
	Right        Node
	IsOverMaxKey bool
}

// GetOverflowKey ...
func (i *InsertResult) GetOverflowKey() []byte {
	return i.Left.GetMaxKey()
}

// GetOverflowPageID ...
func (i *InsertResult) GetOverflowPageID() int64 {
	return i.Left.GetPageID()
}
