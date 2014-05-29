package data

type hashable struct {
	hash Hash256
	raw  []byte
}

func (h *hashable) Hash() Hash256       { return h.hash }
func (h *hashable) Raw() []byte         { return h.raw }
func (h *hashable) SetHash(hash []byte) { copy(h.hash[:], hash[:]) }
func (h *hashable) SetRaw(raw []byte)   { h.raw = make([]byte, len(raw)); copy(h.raw, raw) }
