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


