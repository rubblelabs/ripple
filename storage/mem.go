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
	nodes map[data.Hash256]data.Hashable
	mu    sync.RWMutex
}

func NewMemoryDB(path string) (*MemoryDB, error) {
	mem := &MemoryDB{
		nodes: make(map[data.Hash256]data.Hashable),
	}
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
		node, err := data.NewDecoder(bytes.NewReader(value)).Prefix()
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

func (mem *MemoryDB) Get(hash data.Hash256) (data.Hashable, error) {
	mem.mu.RLock()
	defer mem.mu.RUnlock()
	node, ok := mem.nodes[hash]
	if !ok {
		return nil, ErrNotFound
	}
	return node, nil
}

func (mem *MemoryDB) Stats() string {
	mem.mu.RLock()
	defer mem.mu.RUnlock()
	return fmt.Sprintf("Nodes:%d", len(mem.nodes))
}

func (mem *MemoryDB) Close() {}
