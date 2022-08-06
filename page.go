package bplustree

// Page represent Disk Page.
// TODO: handling Metadata
type Page struct {
	ID   int64
	Data []byte
}

// NewPage Page constructor
func NewPage(id int64, data []byte) *Page {
	return &Page{
		ID:   id,
		Data: data,
	}
}
