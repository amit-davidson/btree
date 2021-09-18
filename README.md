# btree

[![made-with-Go](https://github.com/go-critic/go-critic/workflows/Go/badge.svg)](http://golang.org)
[![made-with-Go](https://img.shields.io/badge/Made%20with-Go-1f425f.svg)](http://golang.org)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://lbesson.mit-license.org/)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)
[![amit-davidson](https://circleci.com/gh/amit-davidson/btree.svg?style=svg)](https://app.circleci.com/pipelines/github/amit-davidson/btree)

btree is a Go implementation of a B-Tree. This project is intended for learning purposes so the code is relatively small
(<500LOC) and highly documented. It means it can be a good starting point for people interested in data structures or 
how databases work.

You can checkout [my blog post](https://medium.com/@amitdavidson234/all-about-b-trees-and-databases-8c0697856189) about the implementation and deep cover about B trees and databases

## Installing

To start using btree, install Go and run `go get`:

```sh
$ go get -u github.com/amit-davidson/btree
```

## Usage
```go
package main

import "fmt"

func main()  {
	minimumItemsInNode := DefaultMinItems
	tree := NewTree(minimumItemsInNode)
	value := "0"
	tree.Put(value, value)

	retVal := tree.Find(value)
	fmt.Printf("Returned value is key:%s value:%s \n", retVal.key, retVal.value)

	tree.Remove(value)

	retVal = tree.Find(value)
	fmt.Print("Returned value is nil")
}
```

## Reading the source code
The best places to start are the operations on the tree:

- `tree.Find()` - Find Returns an item according based on the given key by performing a binary search.


- `tree.Put()` - Put adds a key to the tree. It finds the correct node and the insertion index and adds the item. When performing the
 search, the ancestors are returned as well. This way we can iterate over them to check which nodes were modified and
 rebalance by splitting them accordingly. If the root has too many items, then a new root of a new layer is
 created and the created nodes from the split are added as children.


- `tree.Remove()` - Remove removes a key from the tree. It finds the correct node and the index to remove the item from and removes it.
  When performing the search, the ancestors are returned as well. This way we can iterate over them to check which
  nodes were modified and rebalance by rotating or merging the unbalanced nodes. Rotation is done first and if the
  siblings doesn't have enough items, then merging occurs. If the root is without items after a split, then the root is
  removed and the tree is one level shorter.
