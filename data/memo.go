package data

type MemoItem struct {
	MemoType   VariableLength
	MemoData   VariableLength
	MemoFormat VariableLength
}

type Memo struct {
	Memo MemoItem
}

type Memos []Memo
