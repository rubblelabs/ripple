package ledger

import (
	"fmt"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/data/storage"
	"io"
)

type RadixNode struct {
	Node  data.Hashable
	Depth uint8
}

type RadixMap struct {
	root  data.Hash256
	db    storage.NodeDB
	nodes map[data.Hash256]*RadixNode
	full  bool
}

type WalkFunc func(key data.Hash256, node *RadixNode) error

func NewEmptyRadixMap() *RadixMap {
	return &RadixMap{
		nodes: make(map[data.Hash256]*RadixNode),
	}
}

func NewRadixMap(root data.Hash256, db storage.NodeDB) *RadixMap {
	return &RadixMap{
		root:  root,
		db:    db,
		nodes: make(map[data.Hash256]*RadixNode),
	}
}

func (m *RadixMap) Ledger() *data.Ledger {
	return m.nodes[m.root].Node.(*data.Ledger)
}

func (m *RadixMap) Fill() error {
	if m.full {
		return nil
	}
	if err := m.walk(nil, m.root, 0, true); err != nil {
		return err
	}
	m.full = true
	return nil
}

func (m *RadixMap) Walk(f WalkFunc) error {
	return m.walk(f, m.root, 0, false)
}

func (m *RadixMap) walk(f WalkFunc, key data.Hash256, depth uint8, fill bool) error {
	if key.IsZero() {
		return nil
	}
	var node *RadixNode
	if fill {
		var err error
		node = &RadixNode{
			Depth: depth,
		}
		node.Node, err = m.db.Get(key)
		if err != nil {
			return err
		}
		m.nodes[key] = node
	} else {
		var ok bool
		node, ok = m.nodes[key]
		if !ok {
			return fmt.Errorf("Missing hash: %s", key.String())
		}
		if err := f(key, node); err != nil {
			return err
		}
	}
	inner, ok := node.Node.(*data.InnerNode)
	if !ok {
		return nil
	}
	return inner.Each(func(pos int, child data.Hash256) error {
		return m.walk(f, child, depth+1, fill)
	})
}

func (m *RadixMap) Dump(sequence uint32, w io.Writer) error {
	return m.Walk(func(key data.Hash256, n *RadixNode) error {
		_, err := fmt.Fprintf(w, "%d,%s,%d,%s\n", sequence,
			key.TruncatedString(8), n.Depth, n.Node.GetType())
		return err
	})
}

func (m *RadixMap) Summary(summary map[string]uint64) error {
	return m.Walk(func(key data.Hash256, n *RadixNode) error {
		summary[n.Node.GetType()]++
		return nil
	})
}
