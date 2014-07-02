package data

type hashable struct {
	Hash Hash256 `json:"hash"`
}

func (h *hashable) GetHash() *Hash256 { return &h.Hash }
