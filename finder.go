package bplustree

import (
	"bytes"
)

type FindResultType int

const (
	FindResultTypeUnknown FindResultType = iota
	FindResultTypeMatch
	FindResultTypeFirstGraterThanMatch
)

type FindResult struct {
	Index    int
	Type     FindResultType
	KeyValue KeyValue
}

// Finder represent common logic.
type Finder struct{}

func (f Finder) find(childrens []KeyValue, key []byte) *FindResult {
	if len(childrens) == 0 {
		return &FindResult{
			Type: FindResultTypeUnknown,
		}
	}

	for idx, record := range childrens {
		// TODO: go play link
		compared := bytes.Compare(record.Key, key)
		switch compared {
		// Match
		case 0:
			return &FindResult{
				Index:    idx,
				Type:     FindResultTypeMatch,
				KeyValue: record,
			}
		// FirstGreaterThanMatch
		case 1:
			return &FindResult{
				Index: idx,
				Type:  FindResultTypeFirstGraterThanMatch,
				// KeyValue
			}
		}
	}

	// TODO: last idx
	return &FindResult{
		Index: len(childrens),
		Type:  FindResultTypeFirstGraterThanMatch,
		// KeyValue
	}
}
