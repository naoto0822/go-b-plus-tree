package bplustree

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
