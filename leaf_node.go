package bplustree

import "fmt"

var _ Node = (*LeafNode)(nil)

// LeafNode ...
type LeafNode struct {
	Page Page

	Finder
}

func NewLeafNode(page Page) *LeafNode {
	return &LeafNode{
		Page: page,
	}
}

func (l *LeafNode) GetNodeType() NodeType {
	return NodeTypeLeaf
}

func (l *LeafNode) GetMaxKey() []byte {
	maxKeyValue := l.Page.Records[len(l.Page.Records)-1]
	return maxKeyValue.Key
}

func (l *LeafNode) GetPageID() int64 {
	return l.Page.ID
}

func (l *LeafNode) Get(key []byte) (KeyValue, bool) {
	findResult := l.Finder.find(l.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeMatch:
		return findResult.KeyValue, true
	default:
		return KeyValue{}, false
	}
}

func (l *LeafNode) Insert(key, value []byte) error {
	keyValue := NewKeyValue(key, value)

	findResult := l.Finder.find(l.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeMatch:
		l.Page.UpdateAt(findResult.Index, keyValue)
		return nil

	case FindResultTypeFirstGraterThanMatch:
		l.Page.InsertAt(findResult.Index, keyValue)
		return nil

	case FindResultTypeNoRecord:
		l.Page.InsertAt(0, keyValue)
		return nil

	default:
		return fmt.Errorf("unknown type findResult")
	}
}

func (l *LeafNode) Split() {

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

//func (l *LeafNode) Put(key, value []byte) {
//	keyValue := NewKeyValue(key, value)
//	findResult := l.Finder.find(l.Page.Records, key)
//	switch findResult.Type {
//	case FindResultTypeMatch:
//		l.Page.Update(findResult.Index, keyValue)
//	case FindResultTypeFirstGraterThanMatch:
//		l.Page.Insert(findResult.Index, keyValue)
//	default:
//		// empty
//		l.Page.Insert(0, keyValue)
//	}
//}
