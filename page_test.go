package bplustree

import (
	"testing"
)

func TestPage_SerializeAndDeserialize(t *testing.T) {
	t.Run("Serialize And Deserialize", func(t *testing.T) {
		p := Page{
			ID: 1,
		}
		bytes, err := p.Serialize()
		if err != nil {
			t.Errorf("failed to Page Serialize")
		}

		var pp Page
		err = pp.Deserialize(bytes)
		if err != nil {
			t.Errorf("failed to Page Deserialize")
		}
	})
}
