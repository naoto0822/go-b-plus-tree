package bplustree

import (
	"reflect"
	"testing"
)

func TestBaseNode_Find(t *testing.T) {
	type args struct {
		childrens []KeyValue
		key       []byte
	}

	tests := []struct {
		name string
		args args
		want *FindResult
	}{
		{
			name: "no record",
			args: args{
				childrens: []KeyValue{},
				key:       []byte(`a_key`),
			},
			want: &FindResult{
				Type: FindResultTypeNoRecord,
			},
		},
		{
			name: "match",
			args: args{
				childrens: []KeyValue{
					{Key: []byte(`a_key`), Value: []byte(`a_value`)},
					{Key: []byte(`b_key`), Value: []byte(`b_value`)},
					{Key: []byte(`c_key`), Value: []byte(`c_value`)},
				},
				key: []byte(`b_key`),
			},
			want: &FindResult{
				Type:     FindResultTypeMatch,
				Index:    1,
				KeyValue: KeyValue{Key: []byte(`b_key`), Value: []byte(`b_value`)},
			},
		},
		{
			name: "first grater than match",
			args: args{
				childrens: []KeyValue{
					{Key: []byte(`a_key`), Value: []byte(`a_value`)},
					{Key: []byte(`e_key`), Value: []byte(`e_value`)},
					{Key: []byte(`g_key`), Value: []byte(`g_value`)},
				},
				key: []byte(`f_key`),
			},
			want: &FindResult{
				Type:  FindResultTypeFirstGraterThanMatch,
				Index: 2,
				// KeyValue
			},
		},
		{
			name: "first grater than match, last index",
			args: args{
				childrens: []KeyValue{
					{Key: []byte(`a_key`), Value: []byte(`a_value`)},
					{Key: []byte(`e_key`), Value: []byte(`e_value`)},
					{Key: []byte(`g_key`), Value: []byte(`g_value`)},
				},
				key: []byte(`z_key`),
			},
			want: &FindResult{
				Type:  FindResultTypeFirstGraterThanMatch,
				Index: 2,
				// KeyValue
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			base := BaseNode{}
			got := base.find(tt.args.childrens, tt.args.key)
			if !reflect.DeepEqual(tt.want, got) {
				t.Errorf("%s, want: %+v, got: %+v", tt.name, tt.want, got)
			}
		})
	}
}
