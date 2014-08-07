package memdb

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"github.com/rubblelabs/ripple/data"
	"github.com/rubblelabs/ripple/storage"
	"os"
	"strings"
	"sync"
)

type MemoryDB struct {
	nodes map[data.Hash256]data.Storer
	mu    sync.RWMutex
}

func NewEmptyMemoryDB() *MemoryDB {
	return &MemoryDB{
		nodes: make(map[data.Hash256]data.Storer),
	}
}

func readInput(mem *MemoryDB, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		var nodeid data.Hash256
		if _, err := hex.Decode(nodeid[:], []byte(parts[0])); err != nil {
			return err
		}
		value, err := hex.DecodeString(parts[1])
		if err != nil {
			return err
		}
		node, err := data.ReadPrefix(bytes.NewReader(value), nodeid)
		if err != nil {
			return err
		}
		mem.nodes[nodeid] = node
	}
	return scanner.Err()
}

func NewMemoryDB(paths []string) (*MemoryDB, error) {
	mem := NewEmptyMemoryDB()
	for _, path := range paths {
		if err := readInput(mem, path); err != nil {
			return nil, err
		}
	}
	return mem, nil
}

func (mem *MemoryDB) Get(hash data.Hash256) (data.Storer, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()
	node, ok := mem.nodes[hash]
	if !ok {
		return nil, storage.ErrNotFound
	}
	*node.GetHash() = hash
	return node, nil
}

func (mem *MemoryDB) Insert(item data.Storer) error {
	if item.GetHash().IsZero() {
		return fmt.Errorf("Cannot insert unhashed item")
	}
	mem.mu.Lock()
	mem.nodes[*item.GetHash()] = item
	mem.mu.Unlock()
	return nil
}

func (mem *MemoryDB) Ledger() (*data.LedgerSet, error) {
	return data.NewLedgerSet(32570, 32570), nil
}

func (mem *MemoryDB) Stats() string {
	mem.mu.RLock()
	defer mem.mu.RUnlock()
	return fmt.Sprintf("Nodes:%d", len(mem.nodes))
}

func (mem *MemoryDB) Close() error { return nil }
