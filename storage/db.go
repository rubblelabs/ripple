package storage

import (
	"errors"
	"github.com/donovanhide/ripple/data"
)

var ErrNotFound = errors.New("Not found")

type DB interface {
	Ledger() (*data.LedgerSet, error)
	Get(hash data.Hash256) (data.Storer, error)
	Insert(data.Storer) error
	Stats() string
	Close()
}

type IndexedDB interface {
	DB
	Query(*data.Query) ([]data.Storer, error)
	InsertLookup(string, *LookupItem) error
	GetLookups(string) ([]LookupItem, error)
	GetAccount(uint32) *data.Account
}
