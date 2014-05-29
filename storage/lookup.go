package storage

import (
	"encoding"
	"github.com/donovanhide/ripple/data"
	"reflect"
	"sync"
)

type lookup struct {
	m      map[interface{}]uint32
	r      map[uint32]interface{}
	mu     sync.RWMutex
	db     DB
	insert string
}

type LookupItem struct {
	Id    uint32
	Value interface{}
	Human string
}

func newLookup(get, insert string, db DB, typ reflect.Type) (*lookup, error) {
	items, err := db.GetLookups(get)
	if err != nil {
		return nil, err
	}
	m := make(map[interface{}]uint32)
	r := make(map[uint32]interface{})
	for _, item := range items {
		v := reflect.New(typ).Elem()
		reflect.Copy(v, reflect.ValueOf(item.Value))
		m[v.Interface()] = item.Id
		r[item.Id] = v.Interface()
	}
	return &lookup{
		m:      m,
		r:      r,
		db:     db,
		insert: insert,
	}, nil
}

func (l *lookup) add(v interface{}) (uint32, bool) {
	l.mu.RLock()
	if id, ok := l.m[v]; ok {
		l.mu.RUnlock()
		return id, true
	}
	l.mu.RUnlock()
	// Promote lock
	l.mu.Lock()
	defer l.mu.Unlock()
	if id, ok := l.m[v]; ok {
		return id, true
	}
	id := uint32(len(l.m))
	l.m[v] = id
	l.r[id] = v
	return id, false
}

func (l *lookup) get(n uint32) interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.r[n]
}

func (l *lookup) Lookup(value interface{}) (uint32, error) {
	v := reflect.Indirect(reflect.ValueOf(value))
	id, ok := l.add(v.Interface())
	if ok {
		return id, nil
	}
	human, err := v.Interface().(encoding.TextMarshaler).MarshalText()
	if err != nil {
		return id, err
	}
	item := &LookupItem{
		Id:    id,
		Value: v.Slice(0, v.Len()).Interface(),
		Human: string(human),
	}
	err = l.db.InsertLookup(l.insert, item)
	return id, err
}

type AccountLookup struct {
	*lookup
}

func NewAddressLookup(db DB) (*AccountLookup, error) {
	lookup, err := newLookup("GetAccounts", "InsertAccount", db, reflect.TypeOf(data.Account{}))
	if err != nil {
		return nil, err
	}
	//TODO Move to data
	var accountZero data.Account
	if _, err = lookup.Lookup(&accountZero); err != nil {
		return nil, err
	}
	return &AccountLookup{lookup}, nil
}

func (l *AccountLookup) Get(n uint32) *data.Account {
	if account := l.get(n); account != nil {
		a := account.(data.Account)
		return &a
	}
	return nil
}

type RegularKeyLookup struct {
	*lookup
}

func NewRegularKeyLookup(db DB) (*RegularKeyLookup, error) {
	lookup, err := newLookup("GetRegularKeys", "InsertRegularKey", db, reflect.TypeOf(data.RegularKey{}))
	if err != nil {
		return nil, err
	}
	return &RegularKeyLookup{lookup}, nil
}

func (l *RegularKeyLookup) Get(n uint32) *data.RegularKey {
	if regKey := l.get(n); regKey != nil {
		r := regKey.(data.RegularKey)
		return &r
	}
	return nil
}

type PublicKeyLookup struct {
	*lookup
}

func NewPublicKeyLookup(db DB) (*PublicKeyLookup, error) {
	lookup, err := newLookup("GetPublicKeys", "InsertPublicKey", db, reflect.TypeOf(data.PublicKey{}))
	if err != nil {
		return nil, err
	}
	return &PublicKeyLookup{lookup}, nil
}

func (l *PublicKeyLookup) Get(n uint32) *data.PublicKey {
	if publicKey := l.get(n); publicKey != nil {
		p := publicKey.(data.PublicKey)
		return &p
	}
	return nil
}

type CurrencyLookup struct {
	*lookup
}

func NewCurrencyLookup(db DB) (*CurrencyLookup, error) {
	lookup, err := newLookup("GetCurrencies", "InsertCurrency", db, reflect.TypeOf(data.Currency{}))
	if err != nil {
		return nil, err
	}
	//TODO Move to data
	var xrp data.Currency
	if _, err = lookup.Lookup(&xrp); err != nil {
		return nil, err
	}
	return &CurrencyLookup{lookup}, nil
}

func (l *CurrencyLookup) Get(n uint32) *data.Currency {
	if currency := l.get(n); currency != nil {
		c := currency.(data.Currency)
		return &c
	}
	return nil
}
