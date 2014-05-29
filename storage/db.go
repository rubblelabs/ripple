package storage

import (
	"errors"
	"github.com/donovanhide/ripple/data"
)

var ErrNotFound = errors.New("Not found")

type NodeDB interface {
	Get(hash data.Hash256) (data.Hashable, error)
	Stats() string
	Close()
}

type DB interface {
	Insert(interface{}) error
	Get(*data.Query) ([]interface{}, error)
	InsertLookup(string, *LookupItem) error
	GetLookups(string) ([]LookupItem, error)
	Ledger() (*data.LedgerSet, error)
	GetAccount(uint32) *data.Account
}
