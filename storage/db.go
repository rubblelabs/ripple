package storage

import (
	"errors"
	"github.com/donovanhide/ripple/data"
)

var ErrNotFound = errors.New("Not found")

type DB interface {
	Ledger() (*data.LedgerSet, error)
	Get(hash data.Hash256) (data.Hashable, error)
	Insert(data.Hashable) error
	Stats() string
	Close()
}

type IndexedDB interface {
	DB
	Query(*data.Query) ([]data.Hashable, error)
	InsertLookup(string, *LookupItem) error
	GetLookups(string) ([]LookupItem, error)
	GetAccount(uint32) *data.Account
}
