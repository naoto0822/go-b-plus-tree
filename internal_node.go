package bplustree

import (
	"encoding/binary"
	"fmt"
	"strconv"
)

var _ Node = (*InternalNode)(nil)

type InternalNode struct {
	Page Page

	Finder
}

func NewInternalNode(page Page) *InternalNode {
	return &InternalNode{
		Page: page,
	}
}

func (i *InternalNode) GetNodeType() NodeType {
	return NodeTypeInternal
}

func (i *InternalNode) GetMaxKey() []byte {
	maxKeyValue := i.Page.Records[len(i.Page.Records)-1]
	return maxKeyValue.Key
}

func (i *InternalNode) GetPageID() int64 {
	return i.Page.ID
}

func (i *InternalNode) findChildPageID(key []byte) int64 {
	findResult := i.Finder.find(i.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeMatch:
		return decodePageID(findResult.KeyValue.Value)
	case FindResultTypeFirstGraterThanMatch:
		// TODO
		// if findResult.Index == 0 && i.Page.PrevID != 0 {
		// 	panic("TODO: error handling")
		// }

		// TODO:
		// handle idx?

		kv := i.Page.Records[findResult.Index]
		return decodePageID(kv.Value)
	default:
		panic("TODO: error handling")
	}
}

func (i *InternalNode) Insert(key []byte, pageID int64) error {
	value := encodePageID(pageID)
	keyValue := NewKeyValue(key, value)

	findResult := i.Finder.find(i.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeMatch:
		i.Page.UpdateAt(findResult.Index, keyValue)
		return nil

	case FindResultTypeFirstGraterThanMatch:
		i.Page.InsertAt(findResult.Index, keyValue)
		return nil

	case FindResultTypeNoRecord:
		i.Page.InsertAt(0, keyValue)
		return nil

	default:
		return fmt.Errorf("TODO")
	}
}

func (i *InternalNode) Length() int {
	return len(i.Page.Records)
}

func (i *InternalNode) ByteSize() (int, error) {
	bytes, err := i.Page.Serialize()
	if err != nil {
		return 0, err
	}
	return len(bytes), nil
}

func (i *InternalNode) firstChildPageID() int64 {
	return decodePageID(i.Page.Records[0].Value)
}

func (i *InternalNode) lastChildPage() int64 {
	lastIdx := len(i.Page.Records) - 1
	return decodePageID(i.Page.Records[lastIdx].Value)
}

func encodePageID(src int64) []byte {
	return []byte(strconv.FormatInt(src, 10))
}

func decodePageID(src []byte) int64 {
	// TODO: fuzzy decoding...
	d := binary.BigEndian.Uint64(src)
	return int64(d)
}
