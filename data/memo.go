package data

type Memo struct {
	MemoType VariableLength
	MemoData VariableLength
}

type Memos []Memo
