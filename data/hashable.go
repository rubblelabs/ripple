package data

type hashable struct {
	hash Hash256
}

func (h *hashable) Hash() Hash256       { return h.hash }
func (h *hashable) SetHash(hash []byte) { copy(h.hash[:], hash[:]) }
