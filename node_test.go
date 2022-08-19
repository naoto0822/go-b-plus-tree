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
			},
		},
		{
			name: "over",
			args: args{
				childrens: []KeyValue{
					{Key: []byte(`a_key`), Value: []byte(`a_value`)},
					{Key: []byte(`e_key`), Value: []byte(`e_value`)},
					{Key: []byte(`g_key`), Value: []byte(`g_value`)},
				},
				key: []byte(`z_key`),
			},
			want: &FindResult{
				Type: FindResultTypeOver,
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

func TestBaseNode_RangeScan(t *testing.T) {
	type args struct {
		childrens []KeyValue
		startKey  []byte
		endKey    []byte
	}

	type want struct {
		records   []KeyValue
		isLastHit bool
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "no record",
			args: args{
				childrens: nil,
				startKey:  []byte(`a`),
				endKey:    []byte(`c`),
			},
			want: want{
				records:   nil,
				isLastHit: false,
			},
		},
		{
			name: "a <= (a, b) <= c",
			args: args{
				childrens: []KeyValue{
					{Key: []byte(`a`), Value: []byte(`v_a`)},
					{Key: []byte(`b`), Value: []byte(`v_b`)},
				},
				startKey: []byte(`a`),
				endKey:   []byte(`c`),
			},
			want: want{
				records: []KeyValue{
					{Key: []byte(`a`), Value: []byte(`v_a`)},
					{Key: []byte(`b`), Value: []byte(`v_b`)},
				},
				isLastHit: true,
			},
		},
		{
			name: "a <= (b, c) <= c",
			args: args{
				childrens: []KeyValue{
					{Key: []byte(`b`), Value: []byte(`v_b`)},
					{Key: []byte(`c`), Value: []byte(`v_c`)},
				},
				startKey: []byte(`a`),
				endKey:   []byte(`c`),
			},
			want: want{
				records: []KeyValue{
					{Key: []byte(`b`), Value: []byte(`v_b`)},
					{Key: []byte(`c`), Value: []byte(`v_c`)},
				},
				isLastHit: true,
			},
		},
		{
			name: "a <= (c, d, e) <= g",
			args: args{
				childrens: []KeyValue{
					{Key: []byte(`c`), Value: []byte(`v_c`)},
					{Key: []byte(`d`), Value: []byte(`v_d`)},
					{Key: []byte(`e`), Value: []byte(`v_e`)},
				},
				startKey: []byte(`a`),
				endKey:   []byte(`g`),
			},
			want: want{
				records: []KeyValue{
					{Key: []byte(`c`), Value: []byte(`v_c`)},
					{Key: []byte(`d`), Value: []byte(`v_d`)},
					{Key: []byte(`e`), Value: []byte(`v_e`)},
				},
				isLastHit: true,
			},
		},
		{
			name: "b <= (a, b, c, d, e) <= d",
			args: args{
				childrens: []KeyValue{
					{Key: []byte(`a`), Value: []byte(`v_a`)},
					{Key: []byte(`b`), Value: []byte(`v_b`)},
					{Key: []byte(`c`), Value: []byte(`v_c`)},
					{Key: []byte(`d`), Value: []byte(`v_d`)},
					{Key: []byte(`e`), Value: []byte(`v_e`)},
				},
				startKey: []byte(`b`),
				endKey:   []byte(`d`),
			},
			want: want{
				records: []KeyValue{
					{Key: []byte(`b`), Value: []byte(`v_b`)},
					{Key: []byte(`c`), Value: []byte(`v_c`)},
					{Key: []byte(`d`), Value: []byte(`v_d`)},
				},
				isLastHit: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			baseNode := BaseNode{}
			gotRecord, gotIsNext := baseNode.rangeScan(tt.args.childrens, tt.args.startKey, tt.args.endKey)

			if !reflect.DeepEqual(tt.want.records, gotRecord) {
				t.Errorf("%s, want: %+v, got: %+v", tt.name, tt.want.records, gotRecord)
			}

			if tt.want.isLastHit != gotIsNext {
				t.Errorf("%s, want: %v, got: %v", tt.name, tt.want.isLastHit, gotIsNext)
			}
		})
	}

}
