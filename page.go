package bplustree

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// Page represent Disk Page.
type Page struct {
	ID      int64
	PrevID  int64
	NextID  int64
	Records []KeyValue
}

// NewDefaultPage ...
func NewDefaultPage(id int64) Page {
	return Page{
		ID: id,
	}
}

func (p *Page) Insert(index int, keyValue KeyValue) error {
	pageBytes, err := p.Serialize()
	if err != nil {
		return err
	}
	keyValueBytes, err := keyValue.Serialize()
	if err != nil {
		return err
	}
	// TODO: fuzzy...
	if len(pageBytes)*len(keyValueBytes) > PageSize {
		return fmt.Errorf("page is full")
	}

	p.Records[index] = keyValue
	return nil
}

func (p *Page) Update(index int, keyValue KeyValue) error {
	oldKeyValue := p.Records[index]
	oldBytes, err := oldKeyValue.Serialize()
	if err != nil {
		return err
	}
	pageBytes, err := p.Serialize()
	if err != nil {
		return err
	}
	newBytes, err := keyValue.Serialize()
	if err != nil {
		return err
	}
	// TODO: fuzzy...
	if len(pageBytes)+len(newBytes)-len(oldBytes) > PageSize {
		return fmt.Errorf("page is full")
	}

	// TODO: repack...
	p.Records[index] = keyValue
	return nil
}

func (p Page) Serialize() ([]byte, error) {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	err := encoder.Encode(p)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (p *Page) Deserialize(src []byte) error {
	buf := bytes.NewBuffer(src)
	decoder := gob.NewDecoder(buf)
	err := decoder.Decode(p)
	if err != nil {
		return err
	}
	return nil
}
