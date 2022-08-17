package bplustree

import (
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

const (
	defaultExpire    = 1 * time.Minute
	defaultInterval  = 1 * time.Minute
	bufferPoolKeyFmt = "p_%d"
)

// TODO:
// - gocache -> slice and page_table
// - ClockSweep
// - Frame -> Buffer -> Page

// BufferPool ...
type BufferPool struct {
	pool *gocache.Cache
}

// NewBufferPool ...
func NewBufferPool() *BufferPool {
	pool := gocache.New(defaultExpire, defaultInterval)
	return &BufferPool{
		pool: pool,
	}
}

// Get ...
func (b *BufferPool) Get(pageID int64) (*Page, bool) {
	key := b.getKey(pageID)
	page, found := b.pool.Get(key)
	if !found {
		return nil, false
	}
	return page.(*Page), true
}

// Set ...
func (b *BufferPool) Set(pageID int64, page *Page) {
	key := b.getKey(pageID)
	b.pool.Set(key, page, defaultExpire)
}

func (b *BufferPool) getKey(pageID int64) string {
	return fmt.Sprintf(bufferPoolKeyFmt, pageID)
}
