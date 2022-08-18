package bplustree

// BufferPoolManager ...
type BufferPoolManager struct {
	disk *DiskManager
	pool BufferPool
}

// NewBufferPoolManager ...
func NewBufferPoolManager(disk *DiskManager) *BufferPoolManager {
	pool := NewLruBufferPool()

	return &BufferPoolManager{
		disk: disk,
		pool: pool,
	}
}

// FetchPage ...
func (b *BufferPoolManager) FetchPage(pageID int64) (*Page, error) {
	page, found := b.pool.Get(pageID)
	if found {
		return page, nil
	}

	pageData, err := b.disk.Read(pageID)
	if err != nil {
		return nil, err
	}

	fetchedPage := &Page{}
	err = fetchedPage.Deserialize(pageData)
	if err != nil {
		return nil, err
	}

	b.pool.Set(pageID, fetchedPage)
	return fetchedPage, nil
}

// Flush ...
func (b *BufferPoolManager) Flush(page *Page) error {
	// TODO: need purge
	b.pool.Set(page.ID, page)

	bytes, err := page.Serialize()
	if err != nil {
		return err
	}

	err = b.disk.Write(page.ID, bytes)
	if err != nil {
		return err
	}
	return nil
}

// AllocatePage ...
func (b *BufferPoolManager) AllocatePage(nodeType NodeType) *Page {
	pageID := b.disk.Allocate()
	page := NewDefaultPage(pageID, nodeType)
	// TODO necessary?
	// b.pool.Set(pageID, page)
	return page
}
