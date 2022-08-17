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

type Node interface {
	GetNodeType() NodeType
	GetMaxKey() []byte
	GetPageID() int64
	GetRecords() []KeyValue
}

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
)

type FindResult struct {
	Type     FindResultType
	Index    int64
	KeyValue KeyValue
}

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
				// KeyValue: ,
			}
		}
	}

	// FirstGraterThanMatch
	// last idx
	return &FindResult{
		Type:  FindResultTypeFirstGraterThanMatch,
		Index: int64(len(childrens) - 1), // TODO: later
		// KeyValue: ,
	}
}

type InsertResultType int

const (
	InsertResultTypeUnknown InsertResultType = iota
	InsertResultTypeFit
	InsertResultSplit
)

type InsertResult struct {
	Type           InsertResultType
	OverflowKey    []byte
	OverflowPageID int64
	Left           Node
	Right          Node
}
