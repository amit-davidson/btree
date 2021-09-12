package main

var DefaultMinItems = 128

type Item struct {
	key   string
	value interface{}
}

type Node struct {
	bucket     *Tree
	items      []*Item
	childNodes []*Node
}

type Tree struct {
	root *Node
	minItems int
	maxItems int
}

func newItem(key string, value interface{}) *Item {
	return &Item{
		key:   key,
		value: value,
	}
}

func newTreeWithRoot(root *Node, minItems int) *Tree {
	bucket := &Tree{
		root:    root,
	}
	bucket.root.bucket = bucket
	bucket.minItems = minItems
	bucket.maxItems = minItems * 2
	return bucket
}

func NewTree(minItems int) *Tree {
	return newTreeWithRoot(NewEmptyNode(), minItems)
}

// Put adds a key to the tree. It finds the correct node and the insertion index and adds the item. When performing the
// search, the ancestors are returned as well. This way we can iterate over them to check which nodes were modified and
// rebalance by splitting them accordingly. If the root has too many items, then a new root of a new layer is
// created and the created nodes from the split are added as children.
func (b *Tree) Put(key string, value interface{}) {
	// Find the path to the node where the insertion should happen
	i := newItem(key, value)
	insertionIndex, nodeToInsertIn, ancestorsIndexes := b.findKey(i.key, false)
	// Add item to the leaf node
	nodeToInsertIn.addItem(i, insertionIndex)

	ancestors := b.getNodes(ancestorsIndexes)
	// Rebalance the nodes all the way up. Start From one node before the last and go all the way up. Exclude root.
	for i := len(ancestors) - 2; i >= 0; i-- {
		pnode := ancestors[i]
		node := ancestors[i+1]
		nodeIndex := ancestorsIndexes[i+1]
		if node.isOverPopulated() {
			pnode.split(node, nodeIndex)
		}
	}

	// Handle root
	if b.root.isOverPopulated() {
		newRoot := NewNode(b, []*Item{}, []*Node{b.root})
		newRoot.split(b.root, 0)
		b.root = newRoot
	}
}

// Remove removes a key from the tree. It finds the correct node and the index to remove the item from and removes it.
// When performing the search, the ancestors are returned as well. This way we can iterate over them to check which
// nodes were modified and rebalance by rotating or merging the unbalanced nodes. Rotation is done first. If the
// siblings don't have enough items, then merging occurs. If the root is without items after a split, then the root is
// removed and the tree is one level shorter.
func (b *Tree) Remove(key string) {
	// Find the path to the node where the deletion should happen
	removeItemIndex, nodeToRemoveFrom, ancestorsIndexes := b.findKey(key, true)

	if nodeToRemoveFrom.isLeaf() {
		nodeToRemoveFrom.removeItemFromLeaf(removeItemIndex)
	} else {
		affectedNodes := nodeToRemoveFrom.removeItemFromInternal(removeItemIndex)
		ancestorsIndexes = append(ancestorsIndexes, affectedNodes...)
	}

	ancestors := b.getNodes(ancestorsIndexes)
	// Rebalance the nodes all the way up. Start From one node before the last and go all the way up. Exclude root.
	for i := len(ancestors) - 2; i >= 0; i-- {
		pnode := ancestors[i]
		node := ancestors[i+1]
		if node.isUnderPopulated() {
			pnode.rebalanceRemove(ancestorsIndexes[i+1])
		}
	}
	// If the root has no items after rebalancing
	if len(b.root.items) == 0 && len(b.root.childNodes) > 0 {
		b.root = ancestors[1]
	}
}

// Find Returns an item according based on the given key by performing a binary search.
func (b *Tree) Find(key string) *Item {
	index, containingNode, _ := b.findKey(key, true)
	if index == -1 {
		return nil
	}
	return containingNode.items[index]
}

// findKey finds the node with the key, it's index in the parent's items and a list of its ancestors (not including the
// node itself). The parent's items and key are used later for operations such as searching, adding and removing and list
//of ancestors is used for rebalancing. It's also known as breadcrumbs.
// When the item isn't found, if exact is true, then a falsey answer is returned. If exact is false, then the index
// where the item should have been is returned (Used for insertion)
func (b *Tree) findKey(key string, exact bool) (int, *Node, []int) {
	n := b.root

	// Find the path to the node where the deletion should happen
	ancestorsIndexes := []int{0} // index of root
	for true {
		wasFound, index := n.findKey(key)
		if wasFound {
			return index, n, ancestorsIndexes
		} else {
			if n.isLeaf() {
				if exact {
					return -1, nil, nil
				}
				return index, n, ancestorsIndexes
			}
			nextChild := n.childNodes[index]
			ancestorsIndexes = append(ancestorsIndexes, index)
			n = nextChild
		}
	}
	return -1, nil, nil
}

// getNodes returns a list of nodes based on their indexes (the breadcrumbs) from the root
//           p
//       /       \
//     a          b
//  /     \     /   \
// c       d   e     f
// For [0,1,0] -> p,b,e
func (b *Tree) getNodes(indexes []int) []*Node {
	nodes := []*Node{b.root}
	child := b.root
	for i := 1; i < len(indexes); i++ {
		child = child.childNodes[indexes[i]]
		nodes = append(nodes, child)
	}
	return nodes
}

func NewEmptyNode() *Node {
	return &Node{
		items:      []*Item{},
		childNodes: []*Node{},
	}
}

func NewNode(bucket *Tree, value []*Item, childNodes []*Node) *Node {
	return &Node{
		bucket,
		value,
		childNodes,
	}
}

func isLast(index int, parentNode *Node) bool {
	return index == len(parentNode.items)
}

func isFirst(index int) bool {
	return index == 0
}

func (n *Node) isLeaf() bool {
	return len(n.childNodes) == 0
}

func (n *Node) isOverPopulated() bool {
	return len(n.items) > n.bucket.maxItems
}

func (n *Node) isUnderPopulated() bool {
	return len(n.items) < n.bucket.minItems
}

// findKey iterates all the items and finds the key. If the key is found, then the item is returned. If the key isn't
// found then it means we have to keep searching the tree.
func (n *Node) findKey(key string) (bool, int) {
	for i, existingItem := range n.items {
		if key == existingItem.key {
			return true, i
		}

		if key < existingItem.key {
			return false, i
		}
	}
	return false, len(n.items)
}

// addItem adds an item at a given position. If the item is in the end, then the list is appended. Otherwise, the list
// is shifted and the item is inserted.
func (n *Node) addItem(item *Item, insertionIndex int) int {
	if len(n.items) == insertionIndex { // nil or empty slice or after last element
		n.items = append(n.items, item)
		return insertionIndex
	}

	n.items = append(n.items[:insertionIndex+1], n.items[insertionIndex:]...)
	n.items[insertionIndex] = item
	return insertionIndex
}

// addChild adds a child at a given position. If the child is in the end, then the list is appended. Otherwise, the list
// is shifted and the child is inserted.
func (n *Node) addChild(node *Node, insertionIndex int) {
	if len(n.childNodes) == insertionIndex { // nil or empty slice or after last element
		n.childNodes = append(n.childNodes, node)
	}

	n.childNodes = append(n.childNodes[:insertionIndex+1], n.childNodes[insertionIndex:]...)
	n.childNodes[insertionIndex] = node
}

// split rebalances the tree after adding. After insertion the modified node has to be checked to make sure it
// didn't exceed the maximum number of elements. If it did, then it has to be split and rebalanced. The transformation
// is depicted in the graph below. If it's not a leaf node, then the children has to be moved as well as shown.
// This may leave the parent unbalanced by having too many items so rebalancing has to be checked for all the ancestors.
// 	           n                                        n
//                 3                                       3,6
//	      /        \           ------>       /          |          \
//	   a           modifiedNode            a       modifiedNode     c
//   1,2                 4,5,6,7,8            1,2          4,5         7,8
func (n *Node) split(modifiedNode *Node, insertionIndex int) {
	i := 0
	nodeSize := n.bucket.minItems

	for modifiedNode.isOverPopulated() {
		middleItem := modifiedNode.items[nodeSize]
		var newNode *Node
		if modifiedNode.isLeaf() {
			newNode = NewNode(n.bucket, modifiedNode.items[nodeSize+1:], []*Node{})
			modifiedNode.items = modifiedNode.items[:nodeSize]
		} else {
			newNode = NewNode(n.bucket, modifiedNode.items[nodeSize+1:], modifiedNode.childNodes[i+1:])
			modifiedNode.items = modifiedNode.items[:nodeSize]
			modifiedNode.childNodes = modifiedNode.childNodes[:nodeSize+1]
		}
		n.addItem(middleItem, insertionIndex)
		if len(n.childNodes) == insertionIndex+1 { // If middle of list, then move items forward
			n.childNodes = append(n.childNodes, newNode)
		} else {
			n.childNodes = append(n.childNodes[:insertionIndex+1], n.childNodes[insertionIndex:]...)
			n.childNodes[insertionIndex+1] = newNode
		}

		insertionIndex += 1
		i += 1
		modifiedNode = newNode
	}
}

// rebalanceRemove rebalances the tree after a remove operation. This can be either by rotating to the right, to the
// left or by merging. Firstly, the sibling nodes are checked to see if they have enough items for rebalancing
// (>= minItems+1). If they don't have enough items, then merging with one of the sibling nodes occurs. This may leave
// the parent unbalanced by having too little items so rebalancing has to be checked for all the ancestors.
func (n *Node) rebalanceRemove(unbalancedNodeIndex int) {
	pNode := n
	unbalancedNode := pNode.childNodes[unbalancedNodeIndex]

	// Right rotate
	var leftNode *Node
	if unbalancedNodeIndex != 0 {
		leftNode = pNode.childNodes[unbalancedNodeIndex-1]
		if len(leftNode.items) > n.bucket.minItems {
			rotateRight(leftNode, pNode, unbalancedNode, unbalancedNodeIndex)
			return
		}
	}

	// Left Balance
	var rightNode *Node
	if unbalancedNodeIndex != len(pNode.childNodes)-1 {
		rightNode = pNode.childNodes[unbalancedNodeIndex+1]
		if len(rightNode.items) > n.bucket.minItems {
			rotateLeft(unbalancedNode, pNode, rightNode, unbalancedNodeIndex)
			return
		}
	}

	merge(pNode, unbalancedNodeIndex)
}

func (n *Node) removeItemFromLeaf(index int) {
	n.items = append(n.items[:index], n.items[index+1:]...)
}

func (n *Node) removeItemFromInternal(index int) []int {
	// Take element before inorder (The biggest element from the left branch), put it in the removed index and remove
	// it from the original node.
	//          p
	//       /
	//     ..
	//  /     \
	// ..      a

	affectedNodes := make([]int, 0)
	affectedNodes = append(affectedNodes, index)

	aNode := n.childNodes[index]
	for !aNode.isLeaf() {
		traversingIndex := len(n.childNodes) - 1
		aNode = n.childNodes[traversingIndex]
		affectedNodes = append(affectedNodes, traversingIndex)
	}

	n.items[index] = aNode.items[len(aNode.items)-1]
	aNode.items = aNode.items[:len(aNode.items)-1]
	return affectedNodes
}

func rotateRight(aNode, pNode, bNode *Node, bNodeIndex int) {
	// 	           p                                    p
	//                 4                                    3
	//	      /        \           ------>         /          \
	//	   a           b (unbalanced)            a        b (unbalanced)
	//      1,2,3             5                     1,2            4,5

	// Get last item and remove it
	aNodeItem := aNode.items[len(aNode.items)-1]
	aNode.items = aNode.items[:len(aNode.items)-1]

	// Get item from parent node and assign the aNodeItem item instead
	pNodeItemIndex := bNodeIndex - 1
	if isFirst(bNodeIndex) {
		pNodeItemIndex = 0
	}
	pNodeItem := pNode.items[pNodeItemIndex]
	pNode.items[pNodeItemIndex] = aNodeItem

	// Assign parent item to b and make it first
	bNode.items = append([]*Item{pNodeItem}, bNode.items...)

	// If it's a inner leaf then move children as well.
	if !aNode.isLeaf() {
		childNodeToShift := aNode.childNodes[len(aNode.childNodes)-1]
		aNode.childNodes = aNode.childNodes[:len(aNode.childNodes)-1]
		bNode.childNodes = append([]*Node{childNodeToShift}, bNode.childNodes...)
	}
}

func rotateLeft(aNode, pNode, bNode *Node, bNodeIndex int) {
	// 	           p                                     p
	//                 2                                     3
	//	      /        \           ------>         /          \
	//  a(unbalanced)       b                 a(unbalanced)        b
	//   1                3,4,5                   1,2             4,5

	// Get first item and remove it
	bNodeItem := bNode.items[0]
	bNode.items = bNode.items[1:]

	// Get item from parent node and assign the bNodeItem item instead
	pNodeItemIndex := bNodeIndex
	if isLast(bNodeIndex, pNode) {
		pNodeItemIndex = len(pNode.items) - 1
	}
	pNodeItem := pNode.items[pNodeItemIndex]
	pNode.items[pNodeItemIndex] = bNodeItem

	// Assign parent item to a and make it last
	aNode.items = append(aNode.items, pNodeItem)

	// If it's a inner leaf then move children as well.
	if !bNode.isLeaf() {
		childNodeToShift := bNode.childNodes[0]
		bNode.childNodes = bNode.childNodes[1:]
		aNode.childNodes = append(aNode.childNodes, childNodeToShift)
	}
}

func merge(pNode *Node, unbalancedNodeIndex int) {
	unbalancedNode := pNode.childNodes[unbalancedNodeIndex]
	if unbalancedNodeIndex == 0 {
		// 	               p                                     p
		//                    2,5                                     5
		//	      /        |       \       ------>         /          \
		//  a(unbalanced)   b           c                     a            c
		//   1             3,4          6,7                 1,2,3,4        6,7
		aNode := unbalancedNode
		bNode := pNode.childNodes[unbalancedNodeIndex+1]

		// Take the item from the parent, remove it and add it to the unbalanced node
		pNodeItem := pNode.items[0]
		pNode.items = pNode.items[1:]
		aNode.items = append(aNode.items, pNodeItem)

		//merge the bNode to aNode and remove it. Handle its child nodes as well.
		aNode.items = append(aNode.items, bNode.items...)
		pNode.childNodes = append(pNode.childNodes[0:1], pNode.childNodes[2:]...)
		if !bNode.isLeaf() {
			aNode.childNodes = append(aNode.childNodes, bNode.childNodes...)
		}
	} else {
		// 	               p                                     p
		//                    3,5                                    5
		//	      /        |       \       ------>         /          \
		//           a   b(unbalanced)   c                    a            c
		//          1,2         4        6,7                 1,2,3,4         6,7
		bNode := unbalancedNode
		aNode := pNode.childNodes[unbalancedNodeIndex-1]

		// Take the item from the parent, remove it and add it to the unbalanced node
		pNodeItem := pNode.items[unbalancedNodeIndex-1]
		pNode.items = append(pNode.items[:unbalancedNodeIndex-1], pNode.items[unbalancedNodeIndex:]...)
		aNode.items = append(aNode.items, pNodeItem)

		aNode.items = append(aNode.items, bNode.items...)
		pNode.childNodes = append(pNode.childNodes[:unbalancedNodeIndex], pNode.childNodes[unbalancedNodeIndex+1:]...)
		if !aNode.isLeaf() {
			bNode.childNodes = append(aNode.childNodes, bNode.childNodes...)
		}
	}
}
