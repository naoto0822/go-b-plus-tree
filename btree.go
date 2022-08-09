package bplustree

// BTree ...
type BTree struct {
	bufferPoolManager *BufferPoolManager
}

func NewBTree(bufferPoolManager *BufferPoolManager) *BTree {
	// TODO: root

	return &BTree{
		bufferPoolManager: bufferPoolManager,
	}
}

func (b *BTree) Get() {}

func (b *BTree) Put() {}

func (b *BTree) Delete() {}
