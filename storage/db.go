package storage

import (
	"errors"
	"github.com/rubblelabs/ripple/data"
)

var ErrNotFound = errors.New("Not found")

type DB interface {
	Ledger() (*data.LedgerSet, error)
	Get(hash data.Hash256) (data.Storer, error)
	Insert(data.Storer) error
	Stats() string
	Close() error
}
