package bplustree

import (
	"os"
	"reflect"
	"testing"
)

type testBufferPool struct{}

func (b *testBufferPool) Get(pageID int64) (*Page, bool) {
	return nil, false
}

func (b *testBufferPool) Set(pageID int64, page *Page) {}

func TestBufferPoolManager_FetchPageAndFlush(t *testing.T) {
	t.Run("FetchPage", func(t *testing.T) {
		// Setup
		path := "./test/test_buffer_pool_manager_fetch_page.btr"

		bufferPool := &testBufferPool{}
		diskManager, err := NewDiskManager(path)
		if err != nil {
			t.Errorf("failed to NewDiskManager, err: %v", err)
		}
		bufferPoolManager := &BufferPoolManager{
			disk: diskManager,
			pool: bufferPool,
		}

		// 1. Write
		p1 := &Page{
			ID:       0,
			NodeType: NodeTypeLeaf,
			PrevID:   NoSiblingPageID,
			NextID:   NoSiblingPageID,
			Records:  nil,
		}
		err = bufferPoolManager.Flush(p1)
		if err != nil {
			t.Errorf("failed to Flush, err: %+v", err)
		}

		// 2. Fetch (Exists ID)
		want1 := &Page{
			ID:       0,
			NodeType: NodeTypeLeaf,
			PrevID:   NoSiblingPageID,
			NextID:   NoSiblingPageID,
			Records:  nil,
		}
		got1, err := bufferPoolManager.FetchPage(0)
		if err != nil {
			t.Errorf("failed to FetchPage, err: %+v", err)
		}
		if !reflect.DeepEqual(want1, got1) {
			t.Errorf("FetchPage, want: %+v, got: %+v", want1, got1)
		}

		// 3. Fetch (No Exists ID)
		got2, err := bufferPoolManager.FetchPage(5)
		if err == nil {
			t.Errorf("expect error")
		}
		if got2 != nil {
			t.Errorf("expect not nil, got: %+v", got2)
		}

		// Cleanup
		err = os.Remove(path)
		if err != nil {
			t.Errorf("failed to os.Remove, err: %v", err)
		}
	})
}

func TestBufferPoolManager_AllocatePage(t *testing.T) {
	type args struct {
		disk     *DiskManager
		nodeType NodeType
	}

	tests := []struct {
		name string
		args args
		want *Page
	}{
		{
			name: "already 0 Allocate",
			args: args{
				disk:     &DiskManager{nextPageID: 0},
				nodeType: NodeTypeLeaf,
			},
			want: &Page{
				ID:       0,
				NodeType: NodeTypeLeaf,
				PrevID:   NoSiblingPageID,
				NextID:   NoSiblingPageID,
				Records:  nil,
			},
		},
		{
			name: "already 2 Allocate",
			args: args{
				disk:     &DiskManager{nextPageID: 2},
				nodeType: NodeTypeLeaf,
			},
			want: &Page{
				ID:       2,
				NodeType: NodeTypeLeaf,
				PrevID:   NoSiblingPageID,
				NextID:   NoSiblingPageID,
				Records:  nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufferPool := &testBufferPool{}
			bufferPoolManager := &BufferPoolManager{
				disk: tt.args.disk,
				pool: bufferPool,
			}

			got := bufferPoolManager.AllocatePage(tt.args.nodeType)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("%s, want: %+v, got: %+v", tt.name, tt.want, got)
			}
		})
	}
}
