package bplustree

// KeyValue represent record
type KeyValue struct {
	Key   []byte
	Value []byte
}

// NewKeyValue ...
func NewKeyValue(k, v []byte) KeyValue {
	return KeyValue{
		Key:   k,
		Value: v,
	}
}
