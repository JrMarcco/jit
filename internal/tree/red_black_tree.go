package tree

import (
	"errors"

	"github.com/JrMarcco/easy_kit"
)

var (
	ErrSameRBNode   = errors.New("[easy_kit] cannot insert same red-black tree node")
	ErrNodeNotFound = errors.New("[easy_kit] cannot find node in red-black tree")
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

func (rbt *RBTree[K, V]) Put(key K, val V) error {
	return rbt.insertNode(newNode(key, val))
}

func (rbt *RBTree[K, V]) Del(key K) (V, error) {
	panic("not implemented")
}

// Set sets the value of the given key
func (rbt *RBTree[K, V]) Set(key K, val V) error {
	if node := rbt.findNode(key); node != nil {
		node.val = val
		return nil
	}

	return ErrNodeNotFound
}

func (rbt *RBTree[K, V]) Get(key K) (V, error) {
	if node := rbt.findNode(key); node != nil {
		return node.val, nil
	}

	var zero V
	return zero, ErrNodeNotFound
}

// Keys returns the keys of the tree
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

// Vals returns the values of the tree
func (rbt *RBTree[K, V]) Vals() []V {
	vals := make([]V, 0, rbt.size)

	if rbt.root == nil {
		return vals
	}

	rbt.midOrderTraversal(func(node *rbNode[K, V]) {
		vals = append(vals, node.val)
	})
	return vals
}

// Kvs returns the keys and values of the tree
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

	// the first focus on node is the inserted node
	insertedNode := newNode(node.key, node.val)
	insertedNode.parent = parent

	if cmp < 0 {
		parent.left = insertedNode
	} else {
		parent.right = insertedNode
	}

	rbt.size++
	rbt.fixupInsertion(insertedNode)
	return nil
}

// fixupInsertion ensures the red-black tree properties are maintained after insertion.
// It handles three cases based on the color of the node's uncle:
// 1. Uncle is red
// 2. Uncle is black and its parent is left child
// 3. Uncle is black and its parent is right child
func (rbt *RBTree[K, V]) fixupInsertion(node *rbNode[K, V]) {
	for node != nil && node != rbt.root && node.parent.getColor() == red {
		uncle := node.getUncle()
		if uncle.getColor() == red {
			node = rbt.fixupRedUncle(node, uncle)
			continue
		}

		if node.parent == node.getGrandparent().left {
			// case 2: uncle is black and its parent is left child
			node = rbt.fixupBlackUncleLeftChild(node)
			continue
		}

		// case 3: uncle is black and its parent is right child
		node = rbt.fixupBlackUncleRightChild(node)
	}

	// the new inserted node is root or the new inserted node's parent node is black,
	// no need to fixup
	rbt.root.setColor(black)
}

// fixupRedUncle handles the case where the node's uncle is red.
// It recolors the parent and uncle to black and the grandparent to red,
// then moves the focus to the grandparent.
//
//			    b(c)							           r(c)
//			  /	    \							         /     \
//		  r(b)		r(d)							  b(b)	   b(d)
//		 /	\	    /	\		       =>		     /	\	   /   \
//	  b(z)	r(x)  b(p)  b(q)					  b(z)	r(x)  b(p)  b(q)
//
// 1. change focus on node x's parent node b and uncle node d to black.
func (rbt *RBTree[K, V]) fixupRedUncle(node *rbNode[K, V], uncle *rbNode[K, V]) *rbNode[K, V] {
	grandparent := node.getGrandparent()

	node.parent.setColor(black)
	uncle.setColor(black)
	grandparent.setColor(red)

	return grandparent
}

// fixupBlackUncleLeftChild handles the case where the node's uncle is black,
// and its parent is left child.
//
//			    b(c)							           b(c)
//			  /	    \							         /     \
//		  b(b)		r(d)							  b(b)	   r(d)
//		 /	\	    /  \		       =>		     /	\	   /   \
//	  b(z)	r(y) b(p)  b(q)					      b(z)	r(x) b(p)  b(q)
//			   \							       		/
//		       r(x)								     r(y)
//
// 1. if the focus on node x is right child of its parent, change focus on node to its parent node y.
// 2. left rotate around the focus on node y.
//
//			   b(c)	   	                   b(c)                     r(b)
//			  /	  \	   		              /   \                    /   \
//		  b(b)	   r(d)			      r(b)	   r(d)             b(z)   b(c)
//		 /	\	   /  \	     =>	     /	\	   /  \	     =>            /   \
//	  b(z)	r(x) b(p) b(q)		  b(z)	b(x) b(p) b(q)              b(x)   r(d)
//			/                           /                            /     /  \
//		  r(y)                        r(y)                        r(y)    b(p)  b(q)
//
// 3. change focus on node y's parent to black.
// 4. change focus on node y's grandparent to red.
// 5. right rotate around the focus on node y's grandparent b.
func (rbt *RBTree[K, V]) fixupBlackUncleLeftChild(node *rbNode[K, V]) *rbNode[K, V] {
	if node == node.parent.right {
		node = node.parent
		rbt.leftRotate(node)
	}

	node.parent.setColor(black)
	node.getGrandparent().setColor(red)
	rbt.rightRotate(node.getGrandparent())

	return node.parent
}

// fixupBlackUncleRightChild handles the case where the node's uncle is black,
// and its parent is right child.
//
//			    b(c)							           b(c)
//			  /	    \							         /     \
//		  r(b)		b(d)							  r(b)	   b(d)
//		 /	\	    /  \		       =>		     /	\	   /  \
//	  b(p)	b(q) b(z) r(y)					      b(p)	b(q) b(z) r(x)
//			          /                                             \
//		             r(x)                                           r(y)
//
// 1. if the focus on node x is left child of its parent, change focus on node to its parent node y.
// 2. right rotate around the focus on node y.
//
//			   b(c)	   	                   b(c)                         r(d)
//			  /	  \	   		              /   \                        /   \
//		  r(b)	   b(d)			      r(b)	   r(d)                 b(c)   b(x)
//		 /	\	   /  \	     =>	     /	\	   /  \	     =>         /  \     \
//	  b(p)	b(q) b(z) r(x)		  b(p)	b(q) b(z) b(x)            r(b) b(z)  r(y)
//			          	\                           \             /  \
//		             	r(y)						r(y)	   b(p)  b(q)
//
// 3. change focus on node y's parent to black.
// 4. change focus on node y's grandparent to red.
// 5. left rotate around the focus on node y's grandparent b.
func (rbt *RBTree[K, V]) fixupBlackUncleRightChild(node *rbNode[K, V]) *rbNode[K, V] {
	if node == node.parent.left {
		node = node.parent
		rbt.rightRotate(node)
	}

	node.parent.setColor(black)
	node.getGrandparent().setColor(red)
	rbt.leftRotate(node.getGrandparent())

	return node
}

func (rbt *RBTree[K, V]) deleteNode(node *rbNode[K, V]) (V, error) {
	panic("not implemented")
}

func (rbt *RBTree[K, V]) deletionFixup(node *rbNode[K, V]) {
	panic("not implemented")
}

func (rbt *RBTree[K, V]) findNode(key K) *rbNode[K, V] {
	node := rbt.root
	for node != nil {
		cmp := rbt.cmp(key, node.key)
		if cmp == 0 {
			return node
		}

		if cmp < 0 {
			node = node.left
		} else {
			node = node.right
		}
	}

	return nil
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
