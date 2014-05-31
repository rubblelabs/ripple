package ledger

import (
	"fmt"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/storage"
	"io"
	"strconv"
)

type RadixAction byte

const (
	Addition RadixAction = 'A'
	Deletion RadixAction = 'D'
	Movement RadixAction = 'M'
)

// TODO Replace this with NodeId(Depth!)
type RadixOperation struct {
	*RadixNode
	Action RadixAction
	NodeId data.Hash256
}

func NewRadixOperation(node *RadixNode, action RadixAction, nodeid data.Hash256) *RadixOperation {
	return &RadixOperation{
		RadixNode: node,
		Action:    action,
		NodeId:    nodeid,
	}
}

type RadixDiff struct {
	Left, Right *RadixMap
	Operations  RadixOperations
}

func (r *RadixOperation) String() string {
	return fmt.Sprintf("%c,%d,%s,%s", r.Action, r.Depth, r.NodeId.TruncatedString(8), r.Node.GetType())
}

type RadixOperations []*RadixOperation

func (ro RadixOperations) Len() int { return len(ro) }
func (ro RadixOperations) Less(i, j int) bool {
	if ro[i].Action == ro[j].Action {
		return ro[i].Depth < ro[j].Depth
	}
	return ro[i].Action > ro[j].Action
}
func (ro RadixOperations) Swap(i, j int) { ro[i], ro[j] = ro[j], ro[i] }

func (ro *RadixOperations) Add(node data.Hashable, action RadixAction, nodeid data.Hash256, depth uint8) {
	r := &RadixOperation{
		RadixNode: &RadixNode{
			Node:  node,
			Depth: depth,
		},
		Action: action,
		NodeId: nodeid,
	}
	*ro = append(*ro, r)
}

func (ro RadixOperations) Dump(sequence uint32, w io.Writer) error {
	prefix := strconv.FormatUint(uint64(sequence), 10)
	for _, op := range ro {
		if _, err := fmt.Fprintf(w, "%s,%s\n", prefix, op.String()); err != nil {
			return err
		}
	}
	return nil
}

func Diff(left, right data.Hash256, db storage.DB) (RadixOperations, error) {
	var ops RadixOperations
	if err := diff(left, right, db, &ops, 0); err != nil {
		return nil, err
	}
	return ops, nil
}

func visitChildren(node data.Hashable, db storage.DB, ops *RadixOperations, depth uint8, action RadixAction) error {
	inner, ok := node.(*data.InnerNode)
	if !ok {
		return nil
	}
	return inner.Each(func(pos int, h data.Hash256) error {
		child, err := db.Get(h)
		if err != nil {
			return nil
		}
		ops.Add(child, action, h, depth)
		return visitChildren(child, db, ops, depth+1, action)
	})
}

func diff(left, right data.Hash256, db storage.DB, ops *RadixOperations, depth uint8) error {
	var l, r data.Hashable
	var err error
	switch {
	case left.IsZero() && right.IsZero():
		return nil
	case left.IsZero():
		r, err = db.Get(left)
		if err != nil {
			return err
		}
		ops.Add(r, Addition, right, depth)
		return visitChildren(r, db, ops, depth+1, Deletion)
	case right.IsZero():
		l, err = db.Get(left)
		if err != nil {
			return err
		}
		ops.Add(l, Addition, left, depth)
		return visitChildren(l, db, ops, depth+1, Addition)
	}
	l, err = db.Get(left)
	if err != nil {
		return err
	}
	r, err = db.Get(left)
	if err != nil {
		return err
	}
	ops.Add(r, Addition, right, depth)
	ops.Add(l, Addition, left, depth)
	leftInner, leftOk := l.(*data.InnerNode)
	rightInner, rightOk := r.(*data.InnerNode)
	switch {
	case !leftOk && !rightOk:
		return nil
	case !leftOk:
		return visitChildren(r, db, ops, depth+1, Deletion)
	case !rightOk:
		return visitChildren(l, db, ops, depth+1, Addition)
	default:
		for i := 0; i < 16; i++ {
			leftChild, rightChild := leftInner.Children[i], rightInner.Children[i]
			switch {
			case leftChild == rightChild:
				continue
			case leftChild.IsZero():
				child, err := db.Get(rightChild)
				if err != nil {
					return err
				}
				if err := visitChildren(child, db, ops, depth+1, Deletion); err != nil {
					return err
				}
			case rightChild.IsZero():
				child, err := db.Get(leftChild)
				if err != nil {
					return err
				}
				if err := visitChildren(child, db, ops, depth+1, Addition); err != nil {
					return err
				}
			default:
				if err := diff(leftChild, rightChild, db, ops, depth+1); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (diff *RadixDiff) Dump(sequence uint32, w io.Writer) error {
	for _, op := range diff.Operations {
		_, err := fmt.Fprintf(w, "%d,%c,%s,%d,%s,%d\n", sequence, op.Action,
			op.NodeId.TruncatedString(8), op.Depth, op.Node.GetType())
		if err != nil {
			return err
		}
	}
	return nil
}
