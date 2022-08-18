package bplustree

import (
	"fmt"
	"strconv"
)

var _ Node = (*InternalNode)(nil)

// InternalNode ...
type InternalNode struct {
	Page *Page

	BaseNode
}

// NewInternalNode ...
func NewInternalNode(page *Page) *InternalNode {
	return &InternalNode{
		Page: page,
	}
}

// GetNodeType implements Node
func (i *InternalNode) GetNodeType() NodeType {
	return NodeTypeInternal
}

// GetMaxKey implements Node
func (i *InternalNode) GetMaxKey() []byte {
	maxKeyValue := i.Page.Records[len(i.Page.Records)-1]
	return maxKeyValue.Key
}

// GetPageID implements Node
func (i *InternalNode) GetPageID() int64 {
	return i.Page.ID
}

// GetRecords implements Node
func (i *InternalNode) GetRecords() []KeyValue {
	return i.Page.Records
}

// FindChildPageID ...
func (i *InternalNode) FindChildPageID(key []byte) (int64, bool) {
	findResult := i.BaseNode.find(i.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeNoRecord:
		return 0, false

	case FindResultTypeMatch:
		return decodePageID(findResult.KeyValue.Value), true

	case FindResultTypeFirstGraterThanMatch:
		kv := i.Page.Records[findResult.Index]
		return decodePageID(kv.Value), true

	case FindResultTypeOver:
		kv := i.Page.Records[len(i.Page.Records)-1]
		return decodePageID(kv.Value), true

	default:
		// TODO: error handling
		return 0, false
	}
}

// Insert ...
func (i *InternalNode) Insert(key []byte, pageID int64) error {
	value := encodePageID(pageID)
	keyValue := NewKeyValue(key, value)

	findResult := i.BaseNode.find(i.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeNoRecord:
		i.Page.InsertAt(0, keyValue)
		return nil

	case FindResultTypeMatch:
		i.Page.UpdateAt(findResult.Index, keyValue)
		return nil

	case FindResultTypeFirstGraterThanMatch:
		i.Page.InsertAt(findResult.Index, keyValue)
		return nil

	case FindResultTypeOver:
		index := len(i.Page.Records)
		i.Page.InsertAt(int64(index), keyValue)
		return nil

	default:
		return fmt.Errorf("TODO")
	}
}

// UpdateMaxKey ...
func (i *InternalNode) UpdateMaxKey(key []byte, pageID int64) {
	value := encodePageID(pageID)
	keyValue := NewKeyValue(key, value)
	index := len(i.Page.Records) - 1
	i.Page.UpdateAt(int64(index), keyValue)
}

// Length ...
func (i *InternalNode) Length() int {
	return len(i.Page.Records)
}

func encodePageID(src int64) []byte {
	return []byte(strconv.FormatInt(src, 10))
}

func decodePageID(src []byte) int64 {
	cast := string(src)
	// TODO: error handling
	dest, _ := strconv.ParseInt(cast, 10, 64)
	return dest
}

// String implements Node
func (i *InternalNode) String() string {
	outFmt := "PageID: %d, \n Prev: %d, Next: %d, \n [%s]"
	recordsOut := ""
	for _, r := range i.Page.Records {
		r := fmt.Sprintf("{ %s: %s }", string(r.Key), string(r.Value))
		recordsOut += r
	}
	return fmt.Sprintf(outFmt, i.Page.ID, i.Page.PrevID, i.Page.NextID, recordsOut)
}
