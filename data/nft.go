package data

type NFToken struct {
	NFTokenID *Hash256        `json:",omitempty"`
	URI       *VariableLength `json:",omitempty"`
}