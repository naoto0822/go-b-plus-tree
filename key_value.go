package bplustree

import (
	"bytes"
	"encoding/gob"
)

// KeyValue represent record
type KeyValue struct {
	Key   []byte
	Value []byte
}

func NewKeyValue(k, v []byte) KeyValue {
	return KeyValue{
		Key:   k,
		Value: v,
	}
}

func (k KeyValue) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(k)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (k *KeyValue) Deserialize(src []byte) error {
	buf := bytes.NewBuffer(src)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(k)
	if err != nil {
		return err
	}
	return nil
}
