package tree

import (
	"errors"

	"github.com/JrMarcco/easy_kit"
)

var (
	ErrSameRBNode = errors.New("[easy_kit] cannot insert same red-black tree node")
)

type color int8

const (
	red color = iota
	black
)

// RBTree is a red-black tree
//  1. root node is black
//  2. every leaf node is black or nil, that means leaf node do not store any value
//     2.1 also for space saving implementation share a black empty node
//  3. any neighboring node (parent and child) cannot be red at the same time
//  4. every path from root to leaf node has the same number of black nodes
type RBTree[K any, V any] struct {
	root *rbNode[K, V]
	size int
	cmp  easy_kit.Comparator[K]
}

func NewRBTree[K any, V any](cmp easy_kit.Comparator[K]) *RBTree[K, V] {
	return &RBTree[K, V]{
		root: nil,
		size: 0,
		cmp:  cmp,
	}
}

func (rbt *RBTree[K, V]) Size() int {
	return rbt.size
}

func (rbt *RBTree[K, V]) Insert(key K, val V) error {
	return rbt.insertNode(newNode(key, val))
}

// leftRotate left rotate around the node
//
//	     left rotate around the node x
//		 (a / b / r can be subtree of nil)
//
//		         |                      |
//		         x                      y
//		        / \                    / \
//			   a   y        =>        x   r
//			      / \                / \
//			     b   r              a   b
func (rbt *RBTree[K, V]) leftRotate(x *rbNode[K, V]) {
	if x == nil || x.right == nil {
		// if node x is nil or node x's right is nil, do nothing
		return
	}

	// node y is x's right child
	y := x.right
	// node x's right = node y's left
	x.right = y.left
	// if node y's left is not nil, node y's left's parent = node x
	if y.left != nil {
		y.left.parent = x
	}

	// node y's parent = node x's parent
	y.parent = x.parent

	if x.parent == nil {
		// if node x's parent is nil, node x is root, change root to node y
		rbt.root = y
	} else if x == x.parent.left {
		// if node x is left child, node y is left child
		x.parent.left = y
	} else {
		// if node x is right child, node y is right child
		x.parent.right = y
	}

	// node y's left = node x
	y.left = x
	// node x's parent = node y
	x.parent = y
}

// rightRotate right rotate around the node
//
//	     right rotate around the node x
//		(a / b / r can be subtree of nil)
//
//		     	 |                      |
//		     	 x                      y
//		  	   / \                     / \
//		  	  y   r        =>         a   x
//		  	 / \                    	 / \
//		 	a   b                  		b   r
func (rbt *RBTree[K, V]) rightRotate(x *rbNode[K, V]) {
	if x == nil || x.left == nil {
		// if node x is nil or node x's left is nil, do nothing
		return
	}

	// left: node y
	y := x.left
	// node x's left = node y's right
	x.left = y.right
	// if node y's right is not nil, node y's right's parent = node x
	if y.right != nil {
		y.right.parent = x
	}

	// node y's parent = node x's parent
	y.parent = x.parent

	if x.parent == nil {
		// if node x's parent is nil, node x is root, change root to node y
		rbt.root = y
	} else if x == x.parent.right {
		// if node x is right child, node y is right child
		x.parent.right = y
	} else {
		// if node x is left child, node y is left child
		x.parent.left = y
	}

	// node y's right = node x
	y.right = x
	// node x's parent = node y
	x.parent = y
}

// insertNode insert a new node into the tree
// red-black specifies that the inserted node must be red.
func (rbt *RBTree[K, V]) insertNode(node *rbNode[K, V]) error {

	if rbt.root == nil {
		// if the tree is empty, the inserted node is the root
		rbt.root = newNode(node.key, node.val)
		rbt.root.setColor(black)
		rbt.size++
		return nil
	}

	// the first focus on node is the inserted node
	var focusOnNode *rbNode[K, V]

	cmp := 0
	parent := &rbNode[K, V]{}

	currNode := rbt.root
	for currNode != nil {
		parent = currNode

		cmp = rbt.cmp(node.key, currNode.key)
		if cmp == 0 {
			return ErrSameRBNode
		}

		if cmp < 0 {
			currNode = currNode.left
		} else {
			currNode = currNode.right
		}
	}

	focusOnNode = newNode(node.key, node.val)
	focusOnNode.parent = parent

	if cmp < 0 {
		parent.left = focusOnNode
	} else {
		parent.right = focusOnNode
	}

	rbt.size++
	rbt.fixupInsertion(focusOnNode)
	return nil
}

// fixupInsertion fix the tree after inserting a new node.
// the new node is red before insert fixup.
//
// the focus on node is the inserted node.
//
// case 1: the focus on node's uncle node is red
// case 2: the focus on node's uncle node is black and the focus on node is the right child of its parent
// case 3: the focus on node's uncle node is black and the focus on node is the left child of its parent
func (rbt *RBTree[K, V]) fixupInsertion(node *rbNode[K, V]) {
	uncle := node.getUncle()

	if uncle.getColor() == red {
		// case 1: the focus on node's uncle node is red
		rbt.fixupRedUncle(node, uncle)
	} else {
		// case 2 or case 3: the focus on node's uncle node is black
		rbt.fixupBlackUncle(node)
	}

	// the new inserted node is root or the new inserted node's parent node is black,
	// no need to fixup
	rbt.root.setColor(black)
}

// fixupRedUncle fixup case 1 and return new focus on node
// the focus on node's uncle node is red
//
//			    b(c)							           r(c)
//			  /	    \							         /     \
//		  r(b)		r(d)							  b(b)	   b(d)
//		 /	\	    /	\		       =>		     /	\	   /   \
//	  b(z)	r(a)  b(p)  b(q)					  b(z)	r(a)  b(p)  b(q)
//			/  \							       		/  \
//		b(x)   b(y)								    b(x)   b(y)
//
// focus on node a's uncle node d is red.
// 1. change focus on node's parent node b and uncle node d to black
// 2. change focus on node's grandparent node c to red
// 3. change focus on node to its grandparent node c
// 4. jump to fixupBlackUncle* (case 2 or case 3)
func (rbt *RBTree[K, V]) fixupRedUncle(node *rbNode[K, V], uncle *rbNode[K, V]) {
	grandparent := node.getGrandparent()

	node.parent.setColor(black)
	uncle.setColor(black)
	grandparent.setColor(red)

	rbt.fixupBlackUncle(grandparent)
}

// fixupBlackUncle fixup case 2 or case 3
func (rbt *RBTree[K, V]) fixupBlackUncle(node *rbNode[K, V]) {
	if node.parent == nil {
		return
	}

	if node == node.parent.right {
		// case 2: the focus on node is the right child of its parent
		rbt.fixupBlackUncleRightChild(node)
	} else {
		// case 3: the focus on node is the left child of its parent
		rbt.fixupBlackUncleLeftChild(node)
	}
}

// fixupBlackUncleRightChild fixup case 2
// the focus on node's uncle node is black and the focus on node is the right child of its parent
//
//			    b(c)							           b(c)
//			  /	    \							         /     \
//		  r(b)		b(d)							  r(a)	   b(d)
//		 /	\	    /	\		       =>		     /	\	   /   \
//	  b(z)	r(a)  b(p)  b(q)					  r(b)	b(y)  b(p)  b(q)
//			/  \							     /  \
//		b(x)   b(y)							  b(z)  b(x)
//
// focus on node a's uncle node d is black, and a is the right child of its parent b.
// 1. change focus on node to its parent node（focus on node b）
// 2. left rotate around focus on node b
// 3. jump to fixupBlackUncleAndLeftChild (case 3)
func (rbt *RBTree[K, V]) fixupBlackUncleRightChild(node *rbNode[K, V]) {
	node = node.parent
	rbt.leftRotate(node)

	rbt.fixupBlackUncleLeftChild(node)
}

// fixupBlackUncleLeftChild fixup case 3
// the focus on node's uncle node is black and the focus on node is the left child of its parent
//
//			      b(c)					          r(b)								  b(b)
//				/     \					        /      \							/     \
//	        r(b)	  b(d)				      r(a)	   b(c)						r(a)	  r(c)
//	        /  \      /   \		   =>        /  \	   /   \	      =>       /  \       /  \
//	     r(a)  b(y)  b(p)  b(q)           b(z)  b(x)  b(y)  b(d)			b(z)  b(x)  b(y)  b(d)
//	    /   \                                            	/   \		                      /   \
//	 b(z)   b(x)										  b(p)  b(q)					     b(p)  b(q)
//
// focus on node a's uncle node d is black, and a is the left child of its parent b.
// 1. right rotate around focus on node a's grandparent node c
// 2. swap color between focus on node a's parent node b and a's brother node c
// 3. finish
func (rbt *RBTree[K, V]) fixupBlackUncleLeftChild(node *rbNode[K, V]) {
	grandparent := node.getGrandparent()
	if grandparent == nil {
		return
	}

	node.parent.setColor(black)
	grandparent.setColor(red)
	rbt.rightRotate(grandparent)
}

func (rbt *RBTree[K, V]) Keys() []K {
	keys := make([]K, 0, rbt.size)

	if rbt.root == nil {
		return keys
	}

	rbt.midOrderTraversal(func(node *rbNode[K, V]) {
		keys = append(keys, node.key)
	})
	return keys
}

func (rbt *RBTree[K, V]) Values() []V {
	vals := make([]V, 0, rbt.size)

	if rbt.root == nil {
		return vals
	}

	rbt.midOrderTraversal(func(node *rbNode[K, V]) {
		vals = append(vals, node.val)
	})
	return vals
}

func (rbt *RBTree[K, V]) Kvs() ([]K, []V) {
	keys := make([]K, 0, rbt.size)
	vals := make([]V, 0, rbt.size)

	if rbt.root == nil {
		return keys, vals
	}

	rbt.midOrderTraversal(func(node *rbNode[K, V]) {
		keys = append(keys, node.key)
		vals = append(vals, node.val)
	})

	return keys, vals
}

func (rbt *RBTree[K, V]) midOrderTraversal(visitFn func(node *rbNode[K, V])) {
	stack := make([]*rbNode[K, V], 0, rbt.size)

	currNode := rbt.root
	for currNode != nil || len(stack) > 0 {
		for currNode != nil {
			stack = append(stack, currNode)
			currNode = currNode.left
		}

		currNode = stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		visitFn(currNode)

		currNode = currNode.right
	}
}

type rbNode[K any, V any] struct {
	color  color
	key    K
	val    V
	parent *rbNode[K, V]
	left   *rbNode[K, V]
	right  *rbNode[K, V]
}

// newNode create a new node
// the new node is red before insert fixup
func newNode[K any, V any](key K, val V) *rbNode[K, V] {
	return &rbNode[K, V]{
		key:    key,
		val:    val,
		color:  red,
		parent: nil,
		left:   nil,
		right:  nil,
	}
}

func (rbn *rbNode[K, V]) getColor() color {
	if rbn == nil {
		return black
	}
	return rbn.color
}

func (rbn *rbNode[K, V]) setColor(color color) {
	if rbn == nil {
		return
	}
	rbn.color = color
}

// getGrandparent get the grandparent node
func (rbn *rbNode[K, V]) getGrandparent() *rbNode[K, V] {
	if rbn == nil || rbn.parent == nil {
		return nil
	}
	return rbn.parent.parent
}

// getSibling get the sibling(brother) node
func (rbn *rbNode[K, V]) getSibling() *rbNode[K, V] {
	if rbn == nil || rbn.parent == nil {
		return nil
	}
	if rbn == rbn.parent.left {
		return rbn.parent.right
	}
	return rbn.parent.left
}

// findUncle find the uncle node of the given node
func (rbn *rbNode[K, V]) getUncle() *rbNode[K, V] {
	if rbn == nil {
		return nil
	}

	return rbn.parent.getSibling()
}
