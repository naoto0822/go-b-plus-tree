package bplustree

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

func (l *LeafNode) Get(key []byte) (KeyValue, bool) {
	findResult := l.Finder.find(l.Page.Records, key)
	switch findResult.Type {
	case FindResultTypeMatch:
		return findResult.KeyValue, true
	default:
		return KeyValue{}, false
	}
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
