package bplustree

// BufferPoolManager ...
type BufferPoolManager struct {
	disk *DiskManager
	pool *BufferPool
}

// NewBufferPoolManager ...
func NewBufferPoolManager(disk *DiskManager) *BufferPoolManager {
	pool := NewBufferPool()

	return &BufferPoolManager{
		disk: disk,
		pool: pool,
	}
}

func (b *BufferPoolManager) FetchPage(pageID int64) (Page, error) {
	page, found := b.pool.Get(pageID)
	if found {
		return page, nil
	}

	pageData, err := b.disk.Read(pageID)
	if err != nil {
		return Page{}, err
	}

	fetchedPage := NewPage(pageID, pageData)
	b.pool.Set(pageID, fetchedPage)
	return fetchedPage, nil
}

func (b *BufferPoolManager) AllocatePage() (Page, error) {
	pageID := b.disk.Allocate()
	page := NewDefaultPage(pageID)
	b.pool.Set(pageID, page)
	return page, nil
}

func (b *BufferPoolManager) Commit(page Page) error {
	err := b.disk.Write(page.ID, page.Data)
	if err != nil {
		return err
	}
	return nil
}