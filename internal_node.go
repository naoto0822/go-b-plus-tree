package bplustree

import (
	"encoding/binary"
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

func (i *InternalNode) findChildPageID(key []byte) int64 {
	findResult := i.Finder.find(i.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeMatch:
		return decodePageID(findResult.KeyValue.Value)
	case FindResultTypeFirstGraterThanMatch:
		if findResult.Index == 0 && i.Page.PrevID != 0 {
			panic("TODO: error handling")
		}

		// TODO:
		// handle idx?

		kv := i.Page.Records[findResult.Index]
		return decodePageID(kv.Value)
	default:
		panic("TODO: error handling")
	}
}

func (i *InternalNode) firstChildPageID() int64 {
	return decodePageID(i.Page.Records[0].Value)
}

func (i *InternalNode) lastChildPage() int64 {
	lastIdx := len(i.Page.Records) - 1
	return decodePageID(i.Page.Records[lastIdx].Value)
}

func encodePageID(src int) []byte {
	return []byte(strconv.Itoa(src))
}

func decodePageID(src []byte) int64 {
	// TODO: fuzzy decoding...
	d := binary.BigEndian.Uint64(src)
	return int64(d)
}
