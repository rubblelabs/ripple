package storage

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"github.com/donovanhide/ripple/data"
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

func NewMemoryDB(path string) (*MemoryDB, error) {
	mem := NewEmptyMemoryDB()
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ":")
		var key data.Hash256
		if _, err := hex.Decode(key[:], []byte(parts[0])); err != nil {
			return nil, err
		}
		value, err := hex.DecodeString(parts[1])
		if err != nil {
			return nil, err
		}
		node, err := data.ReadPrefix(bytes.NewReader(value))
		if err != nil {
			return nil, err
		}
		mem.nodes[key] = node
	}
	if scanner.Err() != nil {
		return nil, scanner.Err()
	}
	return mem, nil
}

func (mem *MemoryDB) Get(hash data.Hash256) (data.Storer, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()
	node, ok := mem.nodes[hash]
	if !ok {
		return nil, ErrNotFound
	}
	return node, nil
}

func (mem *MemoryDB) Insert(item data.Storer) error {
	if item.Hash().IsZero() {
		return fmt.Errorf("Cannot insert unhashed item")
	}
	mem.mu.Lock()
	mem.nodes[item.Hash()] = item
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

func (mem *MemoryDB) Close() {}
