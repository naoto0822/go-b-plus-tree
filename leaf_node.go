package bplustree

import (
	"bytes"
)

var _ Node = (*LeafNode)(nil)

// LeafNode ...
type LeafNode struct {
	Page Page
}

func (l *LeafNode) GetNodeType() NodeType {
	return NodeTypeLeaf
}

func (l *LeafNode) IsRootNode() bool {
	return false
}

func (l *LeafNode) IsInternalNode() bool {
	return false
}

func (l *LeafNode) IsLeafNode() bool {
	return true
}

func (l *LeafNode) Get(key []byte) (KeyValue, bool) {
	findResult := l.find(key)
	switch findResult.Type {
	case FindResultTypeMatch:
		return findResult.KeyValue, true
	default:
		return KeyValue{}, false
	}
}

func (l *LeafNode) Put(key, value []byte) {
	keyValue := NewKeyValue(key, value)
	findResult := l.find(key)
	switch findResult.Type {
	case FindResultTypeMatch:
		l.Page.Update(findResult.Index, keyValue)
	case FindResultTypeFirstGraterThanMatch:
		l.Page.Insert(findResult.Index, keyValue)
	default:
		// empty
		l.Page.Insert(0, keyValue)
	}
}

func (l *LeafNode) find(key []byte) *FindResult {
	if len(l.Page.Records) == 0 {
		return &FindResult{
			Type: FindResultTypeUnknown,
		}
	}

	for idx, record := range l.Page.Records {
		compared := bytes.Compare(record.Key, key)
		switch compared {
		// Match
		case 0:
			return &FindResult{
				Index:    idx,
				Type:     FindResultTypeMatch,
				KeyValue: record,
			}
		// FirstGreaterThanMatch
		case 1:
			return &FindResult{
				Index: idx,
				Type:  FindResultTypeFirstGraterThanMatch,
				// KeyValue
			}
		}
	}

	return &FindResult{
		Index: len(l.Page.Records),
		Type:  FindResultTypeFirstGraterThanMatch,
		// KeyValue
	}
}
