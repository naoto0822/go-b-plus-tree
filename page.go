package bplustree

// Page represent Disk Page.
// TODO: handling Metadata
type Page struct {
	ID   int64
	Data []byte // TODO: size PageSize
}

// NewDefaultPage ...
func NewDefaultPage(id int64) Page {
	bytes := make([]byte, PageSize)
	return Page{
		ID:   id,
		Data: bytes,
	}
}

// NewPage Page constructor
func NewPage(id int64, data []byte) Page {
	return Page{
		ID:   id,
		Data: data,
	}
}
