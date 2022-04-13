package data

type Memo struct {
	Memo struct {
		MemoType   VariableLength `json:",omitempty"`
		MemoData   VariableLength `json:",omitempty"`
		MemoFormat VariableLength `json:",omitempty"`
	}
}

type Memos []Memo
