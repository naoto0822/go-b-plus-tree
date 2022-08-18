package bplustree

import (
	"bytes"
	"encoding/gob"
)

// Page represent Disk Page.
type Page struct {
	ID       int64
	NodeType NodeType
	PrevID   int64
	NextID   int64
	Records  []KeyValue
}

// NewDefaultPage ...
func NewDefaultPage(id int64, nodeType NodeType) *Page {
	return &Page{
		ID:       id,
		NodeType: nodeType,
		PrevID:   NoSiblingPageID,
		NextID:   NoSiblingPageID,
	}
}

// InsertAt ...
func (p *Page) InsertAt(index int64, keyValue KeyValue) {
	if int64(len(p.Records)) == index {
		p.Records = append(p.Records, keyValue)
		return
	}

	p.Records = append(p.Records[:index+1], p.Records[index:]...)
	p.Records[index] = keyValue
}

// UpdateAt ...
func (p *Page) UpdateAt(index int64, keyValue KeyValue) {
	p.Records[index] = keyValue
}

// Serialize ...
func (p Page) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Deserialize ...
func (p *Page) Deserialize(src []byte) error {
	buf := bytes.NewBuffer(src)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(p)
	if err != nil {
		return err
	}
	return nil
}
