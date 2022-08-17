package bplustree

import (
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type BufferPool interface {
	Get(pageID int64) (*Page, bool)
	Set(pageID int64, page *Page)
}

var _ BufferPool = (*LruBufferPool)(nil)

const (
	defaultExpire    = 1 * time.Minute
	defaultInterval  = 1 * time.Minute
	bufferPoolKeyFmt = "p_%d"
)

// TODO:
// - gocache -> slice and page_table
// - ClockSweep algorithm
// - Frame -> Buffer -> Page

// LruBufferPool ...
type LruBufferPool struct {
	pool *gocache.Cache
}

// NewBufferPool ...
func NewLruBufferPool() *LruBufferPool {
	pool := gocache.New(defaultExpire, defaultInterval)
	return &LruBufferPool{
		pool: pool,
	}
}

// Get ...
func (b *LruBufferPool) Get(pageID int64) (*Page, bool) {
	key := b.getKey(pageID)
	page, found := b.pool.Get(key)
	if !found {
		return nil, false
	}
	return page.(*Page), true
}

// Set ...
func (b *LruBufferPool) Set(pageID int64, page *Page) {
	key := b.getKey(pageID)
	b.pool.Set(key, page, defaultExpire)
}

func (b *LruBufferPool) getKey(pageID int64) string {
	return fmt.Sprintf(bufferPoolKeyFmt, pageID)
}
