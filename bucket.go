package db

type Item struct {
	key   string
	value interface{}
}

func newItem(key string, value interface{}) *Item {
	return &Item{
		key:   key,
		value: value,
	}
}

type Tree struct {
	root *Node
	minItems int
	maxItems int
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
