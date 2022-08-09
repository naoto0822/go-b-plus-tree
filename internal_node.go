package bplustree

import (
	"encoding/binary"
	"strconv"
)

var _ Node = (*InternalNode)(nil)

type InternalNode struct {
	Page Page

	LeafNode
}

func (i *InternalNode) GetNodeType() NodeType {
	return NodeTypeInternal
}

func (i *InternalNode) IsRootNode() bool {
	return false
}

func (i *InternalNode) IsInternalNode() bool {
	return true
}

func (i *InternalNode) IsLeafNode() bool {
	return false
}

func (i *InternalNode) findChildPageID(key []byte) int {
	findResult := i.find(key)
	switch findResult.Type {
	case FindResultTypeMatch:
		return decodePageID(findResult.KeyValue.Value)
	case FindResultTypeFirstGraterThanMatch:
		if findResult.Index == 0 && i.Page.PrevID != 0 {
			panic("TODO: error handling")
		}
		var index int
		if findResult.Index == 0 {
			index = 0
		} else {
			index = findResult.Index - 1
		}
		kv := i.Page.Records[index]
		return decodePageID(kv.Value)
	default:
		panic("TODO: error handling")
	}
}

func (i *InternalNode) firstChildPage() int {
	return decodePageID(i.Page.Records[0].Value)
}

func (i *InternalNode) lastChildPage() int {
	lastIdx := len(i.Page.Records) - 1
	return decodePageID(i.Page.Records[lastIdx].Value)
}

func (i *InternalNode) addChildNode(pageID int, key []byte) {
	childPageID := encodePageID(pageID)
	i.Put(key, childPageID)
}

func encodePageID(src int) []byte {
	return []byte(strconv.Itoa(src))
}

func decodePageID(src []byte) int {
	d := binary.BigEndian.Uint64(src)
	// TODO: fuzzy...
	return int(d)
}
