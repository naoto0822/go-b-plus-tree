# B+Tree implementation in Go

## Overview

> TODO

## Feature

- [x] Get
- [x] Insert
- [x] RangeScan (ASC)
- [ ] Delete
- [x] Split Node
- [ ] Merge Node

## Design

### Layer

> TODO

### Tree / Node

> TODO

### Page Layout

> TODO

## Usage

```go
path := "./test.btr"
disk, err := bplustree.NewDiskManager(path)
if err != nil {
	panic(err)
}
defer disk.Close()

bufferPoolManager := bplustree.NewBufferPoolManager(disk)
tree := bplustree.NewTree(bufferPoolManager)

// Get
got, err := tree.Get([]byte(`m`))
if err != nil {
	fmt.Printf("error: %+v\n", err)
} else {
	fmt.Printf("got: key: %v, value: %v\n", string(got.Key), string(got.Value))
}

// Insert
err := tree.Insert([]byte(`o`), []byte(`ooo`))
if err != nil {
	fmt.Printf("error: %+v\n", err)
}

// RangeScan
start := []byte(`h`)
end := []byte(`o`)
got, err := tree.RangeScan(start, end)
if err != nil {
	fmt.Printf("error: %+v\n", err)
} else {
	outFmt := ""
	for _, kv := range got {
		outFmt += fmt.Sprintf(" {%s: %s} ", string(kv.Key), string(kv.Value))
	}
	fmt.Printf("got startKey: %v, endKey: %v \n records: %v \n", string(start), string(end), outFmt)
}
```

## TODO

### Slotted Page Layout

> TODO

### Clock-Sweep

> TODO

### Memcomparable Format

- [MyRocks record format](https://github.com/facebook/mysql-5.6/wiki/MyRocks-record-format)
- [pingcap/tidb bytes.go](https://github.com/pingcap/tidb/blob/master/util/codec/bytes.go)

### etc

- [ ] arrange `Tree` reciever name
- [ ] update error message
- [ ] fix error handling
- [ ] delete `NodeTypeRoot`
- [ ] more test code

## Test

```
$ make test
```

## Reference

- [riywo/b-plus-tree](https://github.com/riywo/b-plus-tree)
- [KOBA789/relly](https://github.com/KOBA789/relly)
- [totechite/b_plus_tree](https://github.com/totechite/b_plus_tree)

## License

MIT
