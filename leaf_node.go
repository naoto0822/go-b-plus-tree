package bplustree

import (
	"bytes"
	"fmt"
)

var _ Node = (*LeafNode)(nil)

// LeafNode ...
type LeafNode struct {
	Page *Page

	BaseNode
}

func NewLeafNode(page *Page) *LeafNode {
	return &LeafNode{
		Page: page,
	}
}

func (l *LeafNode) GetNodeType() NodeType {
	return NodeTypeLeaf
}

func (l *LeafNode) GetMaxKey() []byte {
	if len(l.Page.Records) == 0 {
		return nil
	}

	maxKeyValue := l.Page.Records[len(l.Page.Records)-1]
	return maxKeyValue.Key
}

func (l *LeafNode) GetPageID() int64 {
	return l.Page.ID
}

func (l *LeafNode) GetRecords() []KeyValue {
	return l.Page.Records
}

func (l *LeafNode) Get(key []byte) (KeyValue, bool) {
	findResult := l.BaseNode.find(l.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeMatch:
		return findResult.KeyValue, true
	default:
		return KeyValue{}, false
	}
}

func (l *LeafNode) Insert(key, value []byte) error {
	keyValue := NewKeyValue(key, value)

	findResult := l.BaseNode.find(l.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeNoRecord:
		l.Page.InsertAt(0, keyValue)
		return nil

	case FindResultTypeMatch:
		l.Page.UpdateAt(findResult.Index, keyValue)
		return nil

	case FindResultTypeFirstGraterThanMatch:
		l.Page.InsertAt(findResult.Index, keyValue)
		return nil

	case FindResultTypeOver:
		index := len(l.Page.Records)
		l.Page.InsertAt(int64(index), keyValue)
		return nil

	default:
		return fmt.Errorf("unknown type findResult")
	}
}

func (l *LeafNode) Length() int {
	return len(l.Page.Records)
}

func (l *LeafNode) ByteSize() (int, error) {
	bytes, err := l.Page.Serialize()
	if err != nil {
		return 0, err
	}
	return len(bytes), nil
}

func (l *LeafNode) IsOverMaxKey(key []byte) bool {
	maxKey := l.GetMaxKey()
	switch bytes.Compare(maxKey, key) {
	case -1:
		return true
	default:
		return false
	}
}

func (l *LeafNode) String() string {
	outFmt := "PageID: %d, \n Prev: %d, Next: %d, \n [%s]"
	recordsOut := ""
	for _, r := range l.Page.Records {
		r := fmt.Sprintf("{ %s: %s }", string(r.Key), string(r.Value))
		recordsOut += r
	}
	return fmt.Sprintf(outFmt, l.Page.ID, l.Page.PrevID, l.Page.NextID, recordsOut)
}
