package bplustree

type InsertResultType int

const (
	InsertResultTypeUnknown InsertResultType = iota
	InsertResultTypeFit
	InsertResultSplit
)

type InsertResult struct {
	Type           InsertResultType
	OverflowKey    []byte
	OverflowPageID int64
	Left           Node
	Right          Node
}
