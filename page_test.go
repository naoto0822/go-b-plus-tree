package bplustree

import (
	"reflect"
	"testing"
)

func TestPage_NewDefaultPage(t *testing.T) {
	t.Run("NewDefaultPage", func(t *testing.T) {
		want := &Page{
			ID:       1,
			NodeType: NodeTypeInternal,
			PrevID:   NoSiblingPageID,
			NextID:   NoSiblingPageID,
			Records:  nil,
		}

		got := NewDefaultPage(1, NodeTypeInternal)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("failed to NewDefaultPage, want: %+v, got: %+v", want, got)
		}
	})
}

func TestPage_InsertAt(t *testing.T) {
	type args struct {
		index    int64
		keyValue KeyValue
	}

	tests := []struct {
		name  string
		input []KeyValue
		args  args
		want  []KeyValue
	}{
		{
			name: "records 1, index 0",
			input: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
			},
			args: args{
				index:    0,
				keyValue: KeyValue{Key: []byte(`k_a`), Value: []byte(`v_a`)},
			},
			want: []KeyValue{
				{[]byte(`k_a`), []byte(`v_a`)},
				{[]byte(`k_0`), []byte(`v_0`)},
			},
		},
		{
			name: "records 3, index 1",
			input: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
				{[]byte(`k_1`), []byte(`v_1`)},
				{[]byte(`k_2`), []byte(`v_2`)},
			},
			args: args{
				index:    1,
				keyValue: KeyValue{Key: []byte(`k_a`), Value: []byte(`v_a`)},
			},
			want: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
				{[]byte(`k_a`), []byte(`v_a`)},
				{[]byte(`k_1`), []byte(`v_1`)},
				{[]byte(`k_2`), []byte(`v_2`)},
			},
		},
		{
			name: "records 2, index 2",
			input: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
				{[]byte(`k_1`), []byte(`v_1`)},
			},
			args: args{
				index:    2,
				keyValue: KeyValue{Key: []byte(`k_a`), Value: []byte(`v_a`)},
			},
			want: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
				{[]byte(`k_1`), []byte(`v_1`)},
				{[]byte(`k_a`), []byte(`v_a`)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Records: tt.input,
			}
			page.InsertAt(tt.args.index, tt.args.keyValue)
			if !reflect.DeepEqual(tt.want, page.Records) {
				t.Errorf("%s, want: %+v, got: %+v", tt.name, tt.want, page.Records)
			}
		})
	}
}

func TestPage_UpdateAt(t *testing.T) {
	type args struct {
		index    int64
		keyValue KeyValue
	}

	tests := []struct {
		name  string
		input []KeyValue
		args  args
		want  []KeyValue
	}{
		{
			name: "records 1, index 0",
			input: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
			},
			args: args{
				index:    0,
				keyValue: KeyValue{Key: []byte(`k_u`), Value: []byte(`v_u`)},
			},
			want: []KeyValue{
				{[]byte(`k_u`), []byte(`v_u`)},
			},
		},
		{
			name: "records 3, index 1",
			input: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
				{[]byte(`k_1`), []byte(`v_1`)},
				{[]byte(`k_2`), []byte(`v_2`)},
			},
			args: args{
				index:    1,
				keyValue: KeyValue{Key: []byte(`k_u`), Value: []byte(`v_u`)},
			},
			want: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
				{[]byte(`k_u`), []byte(`v_u`)},
				{[]byte(`k_2`), []byte(`v_2`)},
			},
		},
		{
			name: "records 2, index 1",
			input: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
				{[]byte(`k_1`), []byte(`v_1`)},
			},
			args: args{
				index:    1,
				keyValue: KeyValue{Key: []byte(`k_u`), Value: []byte(`v_u`)},
			},
			want: []KeyValue{
				{[]byte(`k_0`), []byte(`v_0`)},
				{[]byte(`k_u`), []byte(`v_u`)},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page := &Page{
				Records: tt.input,
			}
			page.UpdateAt(tt.args.index, tt.args.keyValue)
			if !reflect.DeepEqual(tt.want, page.Records) {
				t.Errorf("%s, want: %+v, got: %+v", tt.name, tt.want, page.Records)
			}
		})
	}
}

func TestPage_SerializeAndDeserialize(t *testing.T) {
	t.Run("Serialize And Deserialize", func(t *testing.T) {
		p := Page{
			ID:       1,
			NodeType: NodeTypeLeaf,
			NextID:   2,
			PrevID:   3,
		}
		bytes, err := p.Serialize()
		if err != nil {
			t.Errorf("failed to Page Serialize")
		}

		pp := &Page{}
		err = pp.Deserialize(bytes)
		if err != nil {
			t.Errorf("failed to Page Deserialize")
		}

		want := &Page{
			ID:       1,
			NodeType: NodeTypeLeaf,
			NextID:   2,
			PrevID:   3,
		}
		if !reflect.DeepEqual(want, pp) {
			t.Errorf("failed to Page Deserailize, want: %+v, got: %+v", want, pp)
		}
	})
}
