package main

import (
	"fmt"

	bplustree "github.com/naoto0822/go-b-plus-tree"
)

func main() {
	fmt.Println("[example: insert_get]")

	path := "./test.btr"
	disk, err := bplustree.NewDiskManager(path)
	if err != nil {
		panic(err)
	}
	defer disk.Close()

	bufferPoolManager := bplustree.NewBufferPoolManager(disk)
	tree := bplustree.NewBTree(bufferPoolManager)

	tree.Insert([]byte(`n`), []byte(`nnn`))
	tree.Insert([]byte(`t`), []byte(`ttt`))
	tree.Insert([]byte(`r`), []byte(`rrr`))
	tree.Insert([]byte(`m`), []byte(`mmm`))
	tree.Insert([]byte(`f`), []byte(`fff`))
	tree.Insert([]byte(`g`), []byte(`ggg`))
	tree.Insert([]byte(`h`), []byte(`hhh`))
	tree.Insert([]byte(`e`), []byte(`eee`))
	tree.Insert([]byte(`i`), []byte(`iii`))
	tree.Insert([]byte(`j`), []byte(`jjj`))
	tree.Insert([]byte(`o`), []byte(`ooo`))
	tree.Insert([]byte(`p`), []byte(`ppp`))
	tree.Insert([]byte(`z`), []byte(`zzz`))

	got, err := tree.Get([]byte(`m`))
	if err != nil {
		fmt.Printf("error: %+v\n", err)
	} else {
		fmt.Printf("got: key: %v, value: %v\n", string(got.Key), string(got.Value))
	}
}
