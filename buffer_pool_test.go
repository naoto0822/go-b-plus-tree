package bplustree

import (
	"reflect"
	"testing"
)

func TestBufferPool_GetSet(t *testing.T) {
	t.Run("Get And Set", func(t *testing.T) {
		pool := NewBufferPool()

		got1, found := pool.Get(1)
		if found != false {
			t.Errorf("failed to Get by 1, want: %v, got: %v", false, found)
		}
		if got1 != nil {
			t.Errorf("failed to Get by 1, want: %v, got: %v", nil, got1)
		}

		page1 := &Page{
			ID: 1,
		}
		pool.Set(1, page1)

		got2, found := pool.Get(1)
		if found != true {
			t.Errorf("failed to Get by 1, want: %v, got: %v", true, found)
		}
		want := &Page{
			ID: 1,
		}
		if !reflect.DeepEqual(got2, want) {
			t.Errorf("failed to Get by 1, want: %v, got: %v", want, got2)
		}
	})
}
