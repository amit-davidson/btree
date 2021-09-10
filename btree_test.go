package db

import (
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

const minItems = 2
const mockNumberOfElements = 10

func (n *Node) addChildNode(child *Node) *Node {
	child.bucket = n.bucket
	n.childNodes = append(n.childNodes, child)
	return n
}

func areTreesEqual(t *testing.T, t1, t2 *Tree) {
	areTreesEqualHelper(t, t1.root, t2.root)
}

func areNodesEqual(t *testing.T, n1, n2 *Node) {
	for i := 0; i < len(n1.items); i++ {
		assert.Equal(t, n1.items[i].key, n2.items[i].key)
		assert.Equal(t, n1.items[i].value, n2.items[i].value)
	}
}

func areTreesEqualHelper(t *testing.T, n1, n2 *Node) {
	require.Equal(t, len(n1.items), len(n2.items))
	require.Equal(t, len(n1.childNodes), len(n2.childNodes))

	areNodesEqual(t, n1, n2)
	// Exit condition: child node -> len(n1.childNodes) == 0
	for i := 0; i < len(n1.childNodes); i++ {
		areTreesEqualHelper(t, n1.childNodes[i], n2.childNodes[i])
	}
}

func createTestMockTree() *Tree {
	root := NewEmptyNode()
	root.addItems("2", "5")

	child0 := NewEmptyNode()
	child0.addItems("0", "1")
	root.addChildNode(child0)

	child1 := NewEmptyNode()
	child1.addItems("3", "4")
	root.addChildNode(child1)

	child2 := NewEmptyNode()
	child2.addItems("6", "7", "8", "9")
	root.addChildNode(child2)

	return newTreeWithRoot(root, minItems)
}

func createTestMockTreeWithout7() *Tree {
	root := NewEmptyNode()
	root.addItems("2", "5")

	child0 := NewEmptyNode()
	child0.addItems("0", "1")
	root.addChildNode(child0)

	child1 := NewEmptyNode()
	child1.addItems("3", "4")
	root.addChildNode(child1)

	child2 := NewEmptyNode()
	child2.addItems("6", "8", "9")
	root.addChildNode(child2)

	return &Tree{root: root}
}

func (n *Node) addItems(keys ...string) *Node {
	for _, key := range keys {
		n.items = append(n.items, newItem(key, key))
	}
	return n
}

func Test_BucketAddSingle(t *testing.T) {
	bucket := NewTree(minItems)
	value := "0"
	bucket.Put(value, value)

	root := NewEmptyNode()
	root.addItems("0")
	expectedbucket := &Tree{root: root}
	areTreesEqual(t, expectedbucket, bucket)
}

func Test_BucketRemoveFromRootSingleElement(t *testing.T) {
	bucket := NewTree(minItems)
	value := "0"
	id := "0"
	bucket.Put(id, value)

	// Tree is balanced
	root := NewEmptyNode()
	root.addItems("0")
	expectedbucket := newTreeWithRoot(root, minItems)
	areTreesEqual(t, expectedbucket, bucket)

	bucket.Remove(id)
	expectedbucketAfterRemoval := NewTree(minItems)
	areTreesEqual(t, expectedbucketAfterRemoval, bucket)
}

func Test_BucketAddMultiple(t *testing.T) {
	bucket := NewTree(minItems)
	numOfElements := mockNumberOfElements
	for i := 0; i < numOfElements; i++ {
		istr := strconv.Itoa(i)
		bucket.Put(istr, istr)
	}

	// Tree is balanced
	areTreesEqual(t, createTestMockTree(), bucket)
}

func Test_BucketAddAndRebalanceSplit(t *testing.T) {
	root := NewEmptyNode()
	root.addItems("4")
	bucket := newTreeWithRoot(root, minItems)

	child0 := NewEmptyNode()
	child0.addItems("0", "1", "2", "3")
	root.addChildNode(child0)

	child1 := NewEmptyNode()
	child1.addItems("5", "6", "7", "8")
	root.addChildNode(child1)

	bucket.Put("9","9")

	expectedroot := NewEmptyNode()
	expectedroot.addItems("4", "7")
	expectedbucket := newTreeWithRoot(expectedroot, minItems)

	expectedchild0 := NewEmptyNode()
	expectedchild0.addItems("0", "1", "2", "3")
	expectedroot.addChildNode(expectedchild0)

	expectedchild1 := NewEmptyNode()
	expectedchild1.addItems("5", "6")
	expectedroot.addChildNode(expectedchild1)

	expectedchild2 := NewEmptyNode()
	expectedchild2.addItems("8", "9")
	expectedroot.addChildNode(expectedchild2)

	// Tree is balanced
	areTreesEqual(t, expectedbucket, bucket)
}

func Test_BucketSplitAndMerge(t *testing.T) {
	root := NewEmptyNode()
	root.addItems("4")
	bucket := newTreeWithRoot(root, minItems)

	child0 := NewEmptyNode()
	child0.addItems("0", "1", "2", "3")
	root.addChildNode(child0)

	child1 := NewEmptyNode()
	child1.addItems("5", "6", "7", "8")
	root.addChildNode(child1)

	bucket.Put("9", "9")

	expectedroot := NewEmptyNode()
	expectedroot.addItems("4", "7")
	expectedbucket := newTreeWithRoot(expectedroot, minItems)

	expectedchild0 := NewEmptyNode()
	expectedchild0.addItems("0", "1", "2", "3")
	expectedroot.addChildNode(expectedchild0)

	expectedchild1 := NewEmptyNode()
	expectedchild1.addItems("5", "6")
	expectedroot.addChildNode(expectedchild1)

	expectedchild2 := NewEmptyNode()
	expectedchild2.addItems("8", "9")
	expectedroot.addChildNode(expectedchild2)

	// Tree is balanced
	areTreesEqual(t, expectedbucket, bucket)

	bucket.Remove("9")

	expectedroot = NewEmptyNode()
	expectedroot.addItems("4")
	expectedbucket = newTreeWithRoot(expectedroot, minItems)

	expectedchild0 = NewEmptyNode()
	expectedchild0.addItems("0", "1", "2", "3")
	expectedroot.addChildNode(expectedchild0)

	expectedchild1 = NewEmptyNode()
	expectedchild1.addItems("5", "6", "7", "8")
	expectedroot.addChildNode(expectedchild1)

	areTreesEqual(t, expectedbucket, bucket)

}

func Test_BucketRemoveFromRootWithoutRebalance(t *testing.T) {
	bucket := NewTree(minItems)
	for i := 0; i < mockNumberOfElements; i++ {
		istr := strconv.Itoa(i)
		bucket.Put(istr, istr)
	}

	// Tree is balanced
	areTreesEqual(t, createTestMockTree(), bucket)

	expectedroot := NewEmptyNode()
	expectedroot.addItems("2", "5")
	expectedTree := newTreeWithRoot(expectedroot, minItems)

	child0 := NewEmptyNode()
	child0.addItems("0", "1")
	expectedroot.addChildNode(child0)

	child1 := NewEmptyNode()
	child1.addItems("3", "4")
	expectedroot.addChildNode(child1)

	child2 := NewEmptyNode()
	child2.addItems("6", "8", "9")
	expectedroot.addChildNode(child2)

	// Remove an element
	bucket.Remove("7")
	areTreesEqual(t, expectedTree, bucket)
}

func Test_BucketRemoveFromRootAndRotateLeft(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("2", "5")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("0", "1")
	mockRoot.addChildNode(mockChild0)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("3", "4")
	mockRoot.addChildNode(mockChild1)

	mockChild2 := NewEmptyNode()
	mockChild2.addItems("6", "7", "8")
	mockRoot.addChildNode(mockChild2)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("2", "6")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild0 := NewEmptyNode()
	expectedChild0.addItems("0", "1")
	expectedRoot.addChildNode(expectedChild0)

	expectedChild1 := NewEmptyNode()
	expectedChild1.addItems("3", "4")
	expectedRoot.addChildNode(expectedChild1)

	expectedChild2 := NewEmptyNode()
	expectedChild2.addItems("7", "8")
	expectedRoot.addChildNode(expectedChild2)

	// Remove an element
	mockTree.Remove("5")
	areTreesEqual(t, expectedTree, mockTree)
}

func Test_BucketRemoveFromRootAndRotateRight(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("3", "6")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("0", "1", "2")
	mockRoot.addChildNode(mockChild0)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("4", "5")
	mockRoot.addChildNode(mockChild1)

	mockChild2 := NewEmptyNode()
	mockChild2.addItems("7", "8")
	mockRoot.addChildNode(mockChild2)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("2", "5")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild0 := NewEmptyNode()
	expectedChild0.addItems("0", "1")
	expectedRoot.addChildNode(expectedChild0)

	expectedChild1 := NewEmptyNode()
	expectedChild1.addItems("3", "4")
	expectedRoot.addChildNode(expectedChild1)

	expectedChild2 := NewEmptyNode()
	expectedChild2.addItems("7", "8")
	expectedRoot.addChildNode(expectedChild2)

	// Remove an element
	mockTree.Remove("6")
	areTreesEqual(t, expectedTree, mockTree)
}

// Test_BucketRemoveFromRootAndRebalanceMergeToUnbalanced tests when the unbalanced node is the most left one so the
// merge has to happen from the right node into the unbalanced node
func Test_BucketRemoveFromRootAndRebalanceMergeToUnbalanced(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("2", "5")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("0", "1")
	mockRoot.addChildNode(mockChild0)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("3", "4")
	mockRoot.addChildNode(mockChild1)

	mockChild2 := NewEmptyNode()
	mockChild2.addItems("6", "7")
	mockRoot.addChildNode(mockChild2)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("5")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild0 := NewEmptyNode()
	expectedChild0.addItems("0", "1", "3", "4")
	expectedRoot.addChildNode(expectedChild0)

	expectedChild1 := NewEmptyNode()
	expectedChild1.addItems("6", "7")
	expectedRoot.addChildNode(expectedChild1)

	// Remove an element
	mockTree.Remove("2")
	areTreesEqual(t, expectedTree, mockTree)
}

// Test_BucketRemoveFromRootAndRebalanceMergeFromUnbalanced tests when the unbalanced node is not the most left one so the
// merge has to happen from the unbalanced node to the node left to it
func Test_BucketRemoveFromRootAndRebalanceMergeFromUnbalanced(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("2", "5")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("0", "1")
	mockRoot.addChildNode(mockChild0)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("3", "4")
	mockRoot.addChildNode(mockChild1)

	mockChild2 := NewEmptyNode()
	mockChild2.addItems("6", "7")
	mockRoot.addChildNode(mockChild2)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("4")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild0 := NewEmptyNode()
	expectedChild0.addItems("0", "1", "2", "3")
	expectedRoot.addChildNode(expectedChild0)

	expectedChild1 := NewEmptyNode()
	expectedChild1.addItems("6", "7")
	expectedRoot.addChildNode(expectedChild1)


	// Remove an element
	mockTree.Remove("5")
	areTreesEqual(t, expectedTree, mockTree)
}

func Test_BucketRemoveFromInnerNodeAndRotateLeft(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("8")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("2", "5")
	mockRoot.addChildNode(mockChild0)
	mockChild00 := NewEmptyNode()
	mockChild00.addItems("0", "1")
	mockChild0.addChildNode(mockChild00)
	mockChild01 := NewEmptyNode()
	mockChild01.addItems("3", "4")
	mockChild0.addChildNode(mockChild01)
	mockChild02 := NewEmptyNode()
	mockChild02.addItems("6", "7")
	mockChild0.addChildNode(mockChild02)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("b", "e", "h")
	mockRoot.addChildNode(mockChild1)
	mockChild10 := NewEmptyNode()
	mockChild10.addItems("9", "a")
	mockChild1.addChildNode(mockChild10)
	mockChild11 := NewEmptyNode()
	mockChild11.addItems("c", "d")
	mockChild1.addChildNode(mockChild11)
	mockChild12 := NewEmptyNode()
	mockChild12.addItems("f", "g")
	mockChild1.addChildNode(mockChild12)
	mockChild13 := NewEmptyNode()
	mockChild13.addItems("i", "j")
	mockChild1.addChildNode(mockChild13)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("b")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild0 := NewEmptyNode()
	expectedChild0.addItems("4", "8")
	expectedRoot.addChildNode(expectedChild0)
	expectedChild00 := NewEmptyNode()
	expectedChild00.addItems("0", "1", "2", "3")
	expectedChild0.addChildNode(expectedChild00)
	expectedChild01 := NewEmptyNode()
	expectedChild01.addItems("6", "7")
	expectedChild0.addChildNode(expectedChild01)
	expectedChild02 := NewEmptyNode()
	expectedChild02.addItems("9", "a")
	expectedChild0.addChildNode(expectedChild02)

	expectedChild1 := NewEmptyNode()
	expectedChild1.addItems("e", "h")
	expectedRoot.addChildNode(expectedChild1)
	expectedChild10 := NewEmptyNode()
	expectedChild10.addItems("c", "d")
	expectedChild1.addChildNode(expectedChild10)
	expectedChild11 := NewEmptyNode()
	expectedChild11.addItems("f", "g")
	expectedChild1.addChildNode(expectedChild11)
	expectedChild12 := NewEmptyNode()
	expectedChild12.addItems("i", "j")
	expectedChild1.addChildNode(expectedChild12)

	// Remove an element
	mockTree.Remove("5")
	areTreesEqual(t, expectedTree, mockTree)
}

func Test_BucketRemoveFromInnerNodeAndRotateRight(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("b")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("2", "5", "8")
	mockRoot.addChildNode(mockChild0)
	mockChild00 := NewEmptyNode()
	mockChild00.addItems("0", "1")
	mockChild0.addChildNode(mockChild00)
	mockChild01 := NewEmptyNode()
	mockChild01.addItems("3", "4")
	mockChild0.addChildNode(mockChild01)
	mockChild02 := NewEmptyNode()
	mockChild02.addItems("6", "7")
	mockChild0.addChildNode(mockChild02)
	mockChild03 := NewEmptyNode()
	mockChild03.addItems("9", "a")
	mockChild0.addChildNode(mockChild03)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("e", "h")
	mockRoot.addChildNode(mockChild1)
	mockChild10 := NewEmptyNode()
	mockChild10.addItems("c", "d")
	mockChild1.addChildNode(mockChild10)
	mockChild11 := NewEmptyNode()
	mockChild11.addItems("f", "g")
	mockChild1.addChildNode(mockChild11)
	mockChild12 := NewEmptyNode()
	mockChild12.addItems("i", "j")
	mockChild1.addChildNode(mockChild12)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("8")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild0 := NewEmptyNode()
	expectedChild0.addItems("2", "5")
	expectedRoot.addChildNode(expectedChild0)
	expectedChild00 := NewEmptyNode()
	expectedChild00.addItems("0", "1")
	expectedChild0.addChildNode(expectedChild00)
	expectedChild01 := NewEmptyNode()
	expectedChild01.addItems("3", "4")
	expectedChild0.addChildNode(expectedChild01)
	expectedChild02 := NewEmptyNode()
	expectedChild02.addItems("6", "7")
	expectedChild0.addChildNode(expectedChild02)

	expectedChild1 := NewEmptyNode()
	expectedChild1.addItems("b", "h")
	expectedRoot.addChildNode(expectedChild1)
	expectedChild10 := NewEmptyNode()
	expectedChild10.addItems("9", "a")
	expectedChild1.addChildNode(expectedChild10)
	expectedChild11 := NewEmptyNode()
	expectedChild11.addItems("c", "d", "f", "g")
	expectedChild1.addChildNode(expectedChild11)
	expectedChild12 := NewEmptyNode()
	expectedChild12.addItems("i", "j")
	expectedChild1.addChildNode(expectedChild12)

	// Remove an element
	mockTree.Remove("e")
	areTreesEqual(t, expectedTree, mockTree)
}

func Test_BucketRemoveFromInnerNodeAndUnion(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("8")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("2", "5")
	mockRoot.addChildNode(mockChild0)
	mockChild00 := NewEmptyNode()
	mockChild00.addItems("0", "1")
	mockChild0.addChildNode(mockChild00)
	mockChild01 := NewEmptyNode()
	mockChild01.addItems("3", "4")
	mockChild0.addChildNode(mockChild01)
	mockChild02 := NewEmptyNode()
	mockChild02.addItems("6", "7")
	mockChild0.addChildNode(mockChild02)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("b", "e")
	mockRoot.addChildNode(mockChild1)
	mockChild10 := NewEmptyNode()
	mockChild10.addItems("9", "a")
	mockChild1.addChildNode(mockChild10)
	mockChild11 := NewEmptyNode()
	mockChild11.addItems("c", "d")
	mockChild1.addChildNode(mockChild11)
	mockChild12 := NewEmptyNode()
	mockChild12.addItems("f", "g")
	mockChild1.addChildNode(mockChild12)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("5", "8", "b", "e")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild0 := NewEmptyNode()
	expectedChild0.addItems("0", "1", "3", "4")
	expectedRoot.addChildNode(expectedChild0)
	expectedChild1 := NewEmptyNode()
	expectedChild1.addItems("6", "7")
	expectedRoot.addChildNode(expectedChild1)
	expectedChild2 := NewEmptyNode()
	expectedChild2.addItems("9", "a")
	expectedRoot.addChildNode(expectedChild2)
	expectedChild3 := NewEmptyNode()
	expectedChild3.addItems("c", "d")
	expectedRoot.addChildNode(expectedChild3)
	expectedChild4 := NewEmptyNode()
	expectedChild4.addItems("f", "g")
	expectedRoot.addChildNode(expectedChild4)

	// Remove an element
	mockTree.Remove("2")
	areTreesEqual(t, expectedTree, mockTree)
}

func Test_BucketRemoveFromLeafAndRotateLeft(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("9")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("2", "6")
	mockRoot.addChildNode(mockChild0)
	mockChild00 := NewEmptyNode()
	mockChild00.addItems("0", "1")
	mockChild0.addChildNode(mockChild00)
	mockChild01 := NewEmptyNode()
	mockChild01.addItems("3", "4", "5")
	mockChild0.addChildNode(mockChild01)
	mockChild02 := NewEmptyNode()
	mockChild02.addItems("7", "8")
	mockChild0.addChildNode(mockChild02)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("c", "f")
	mockRoot.addChildNode(mockChild1)
	mockChild10 := NewEmptyNode()
	mockChild10.addItems("a", "b")
	mockChild1.addChildNode(mockChild10)
	mockChild11 := NewEmptyNode()
	mockChild11.addItems("d", "e")
	mockChild1.addChildNode(mockChild11)
	mockChild12 := NewEmptyNode()
	mockChild12.addItems("g", "h")
	mockChild1.addChildNode(mockChild12)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("9")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild0 := NewEmptyNode()
	expectedChild0.addItems("3", "6")
	expectedRoot.addChildNode(expectedChild0)
	expectedChild00 := NewEmptyNode()
	expectedChild00.addItems("0", "2")
	expectedChild0.addChildNode(expectedChild00)
	expectedChild01 := NewEmptyNode()
	expectedChild01.addItems("4", "5")
	expectedChild0.addChildNode(expectedChild01)
	expectedChild02 := NewEmptyNode()
	expectedChild02.addItems("7", "8")
	expectedChild0.addChildNode(expectedChild02)

	expectedChild1 := NewEmptyNode()
	expectedChild1.addItems("c", "f")
	expectedRoot.addChildNode(expectedChild1)
	expectedChild10 := NewEmptyNode()
	expectedChild10.addItems("a", "b")
	expectedChild1.addChildNode(expectedChild10)
	expectedChild11 := NewEmptyNode()
	expectedChild11.addItems("d", "e")
	expectedChild1.addChildNode(expectedChild11)
	expectedChild12 := NewEmptyNode()
	expectedChild12.addItems("g", "h")
	expectedChild1.addChildNode(expectedChild12)

	// Remove an element
	mockTree.Remove("1")
	areTreesEqual(t, expectedTree, mockTree)
}

func Test_BucketRemoveFromLeafAndRotateRight(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("9")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("2", "6")
	mockRoot.addChildNode(mockChild0)
	mockChild00 := NewEmptyNode()
	mockChild00.addItems("0", "1")
	mockChild0.addChildNode(mockChild00)
	mockChild01 := NewEmptyNode()
	mockChild01.addItems("3", "4", "5")
	mockChild0.addChildNode(mockChild01)
	mockChild02 := NewEmptyNode()
	mockChild02.addItems("7", "8")
	mockChild0.addChildNode(mockChild02)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("c", "f")
	mockRoot.addChildNode(mockChild1)
	mockChild10 := NewEmptyNode()
	mockChild10.addItems("a", "b")
	mockChild1.addChildNode(mockChild10)
	mockChild11 := NewEmptyNode()
	mockChild11.addItems("d", "e")
	mockChild1.addChildNode(mockChild11)
	mockChild12 := NewEmptyNode()
	mockChild12.addItems("g", "h")
	mockChild1.addChildNode(mockChild12)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("9")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild0 := NewEmptyNode()
	expectedChild0.addItems("2", "5")
	expectedRoot.addChildNode(expectedChild0)
	expectedChild00 := NewEmptyNode()
	expectedChild00.addItems("0", "1")
	expectedChild0.addChildNode(expectedChild00)
	expectedChild01 := NewEmptyNode()
	expectedChild01.addItems("3", "4")
	expectedChild0.addChildNode(expectedChild01)
	expectedChild02 := NewEmptyNode()
	expectedChild02.addItems("6", "7")
	expectedChild0.addChildNode(expectedChild02)

	expectedChild1 := NewEmptyNode()
	expectedChild1.addItems("c", "f")
	expectedRoot.addChildNode(expectedChild1)
	expectedChild10 := NewEmptyNode()
	expectedChild10.addItems("a", "b")
	expectedChild1.addChildNode(expectedChild10)
	expectedChild11 := NewEmptyNode()
	expectedChild11.addItems("d", "e")
	expectedChild1.addChildNode(expectedChild11)
	expectedChild12 := NewEmptyNode()
	expectedChild12.addItems("g", "h")
	expectedChild1.addChildNode(expectedChild12)

	// Remove an element
	mockTree.Remove("8")
	areTreesEqual(t, expectedTree, mockTree)
}

func Test_BucketRemoveFromLeafAndUnion(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("8")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("2", "5")
	mockRoot.addChildNode(mockChild0)
	mockChild00 := NewEmptyNode()
	mockChild00.addItems("0", "1")
	mockChild0.addChildNode(mockChild00)
	mockChild01 := NewEmptyNode()
	mockChild01.addItems("3", "4")
	mockChild0.addChildNode(mockChild01)
	mockChild02 := NewEmptyNode()
	mockChild02.addItems("6", "7")
	mockChild0.addChildNode(mockChild02)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("b", "e")
	mockRoot.addChildNode(mockChild1)
	mockChild10 := NewEmptyNode()
	mockChild10.addItems("9", "a")
	mockChild1.addChildNode(mockChild10)
	mockChild11 := NewEmptyNode()
	mockChild11.addItems("c", "d")
	mockChild1.addChildNode(mockChild11)
	mockChild12 := NewEmptyNode()
	mockChild12.addItems("f", "g")
	mockChild1.addChildNode(mockChild12)

	expectedRoot := NewEmptyNode()
	expectedRoot.addItems("5", "8", "b", "e")
	expectedTree := newTreeWithRoot(expectedRoot, minItems)

	expectedChild00 := NewEmptyNode()
	expectedChild00.addItems("1", "2", "3", "4")
	expectedRoot.addChildNode(expectedChild00)
	expectedChild01 := NewEmptyNode()
	expectedChild01.addItems("6", "7")
	expectedRoot.addChildNode(expectedChild01)
	expectedChild02 := NewEmptyNode()
	expectedChild02.addItems("9", "a")
	expectedRoot.addChildNode(expectedChild02)
	expectedChild03 := NewEmptyNode()
	expectedChild03.addItems("c", "d")
	expectedRoot.addChildNode(expectedChild03)
	expectedChild04 := NewEmptyNode()
	expectedChild04.addItems("f", "g")
	expectedRoot.addChildNode(expectedChild04)

	// Remove an element
	mockTree.Remove("0")
	areTreesEqual(t, expectedTree, mockTree)
}

func Test_BucketFindNode(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("8")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("2", "5")
	mockRoot.addChildNode(mockChild0)
	mockChild00 := NewEmptyNode()
	mockChild00.addItems("0", "1")
	mockChild0.addChildNode(mockChild00)
	mockChild01 := NewEmptyNode()
	mockChild01.addItems("3", "4")
	mockChild0.addChildNode(mockChild01)
	mockChild02 := NewEmptyNode()
	mockChild02.addItems("6", "7")
	mockChild0.addChildNode(mockChild02)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("b", "e")
	mockRoot.addChildNode(mockChild1)
	mockChild10 := NewEmptyNode()
	mockChild10.addItems("9", "a")
	mockChild1.addChildNode(mockChild10)
	mockChild11 := NewEmptyNode()
	mockChild11.addItems("c", "d")
	mockChild1.addChildNode(mockChild11)
	mockChild12 := NewEmptyNode()
	mockChild12.addItems("f", "g")
	mockChild1.addChildNode(mockChild12)

	// Item found
	expectedItem := newItem("c", "c")
	item := mockTree.Find("c")
	assert.Equal(t, expectedItem, item)

	// Item not found
	expectedItem = nil
	item = mockTree.Find("h")
	assert.Equal(t, expectedItem, item)
}

func Test_BucketUpdateNode(t *testing.T) {
	mockRoot := NewEmptyNode()
	mockRoot.addItems("8")
	mockTree := newTreeWithRoot(mockRoot, minItems)

	mockChild0 := NewEmptyNode()
	mockChild0.addItems("2", "5")
	mockRoot.addChildNode(mockChild0)
	mockChild00 := NewEmptyNode()
	mockChild00.addItems("0", "1")
	mockChild0.addChildNode(mockChild00)
	mockChild01 := NewEmptyNode()
	mockChild01.addItems("3", "4")
	mockChild0.addChildNode(mockChild01)
	mockChild02 := NewEmptyNode()
	mockChild02.addItems("6", "7")
	mockChild0.addChildNode(mockChild02)

	mockChild1 := NewEmptyNode()
	mockChild1.addItems("b", "e")
	mockRoot.addChildNode(mockChild1)
	mockChild10 := NewEmptyNode()
	mockChild10.addItems("9", "a")
	mockChild1.addChildNode(mockChild10)
	mockChild11 := NewEmptyNode()
	mockChild11.addItems("c", "d")
	mockChild1.addChildNode(mockChild11)
	mockChild12 := NewEmptyNode()
	mockChild12.addItems("f", "g")
	mockChild1.addChildNode(mockChild12)

	// Item found
	expectedItem := newItem("c", "c")
	item := mockTree.Find("c")
	assert.Equal(t, expectedItem, item)

	// Item updated successfully
	newvalue := "f"
	mockTree.Put("c",newvalue)
	item = mockTree.Find("c")
	assert.Equal(t, newvalue, item.value)
}
