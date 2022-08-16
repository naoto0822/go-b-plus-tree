package bplustree

import (
	"bytes"
)

type FindResultType int

const (
	FindResultTypeUnknown FindResultType = iota
	FindResultTypeNoRecord
	FindResultTypeMatch
	FindResultTypeFirstGraterThanMatch
)

type FindResult struct {
	Type     FindResultType
	Index    int64
	KeyValue KeyValue
}

// Finder represent common logic (like extend).
type Finder struct{}

func (f Finder) find(childrens []KeyValue, key []byte) *FindResult {
	if len(childrens) == 0 {
		return &FindResult{
			Type: FindResultTypeNoRecord,
		}
	}

	for idx, record := range childrens {
		// TODO: go play link
		compared := bytes.Compare(record.Key, key)
		switch compared {
		// Match
		case 0:
			return &FindResult{
				Type:     FindResultTypeMatch,
				Index:    int64(idx),
				KeyValue: record,
			}

		// FirstGreaterThanMatch
		case 1:
			return &FindResult{
				Type:  FindResultTypeFirstGraterThanMatch,
				Index: int64(idx),
				// KeyValue: ,
			}
		}
	}

	// TODO: last idx
	return &FindResult{
		Type:  FindResultTypeFirstGraterThanMatch,
		Index: int64(len(childrens)),
		// KeyValue: ,
	}
}
