package data

type Memo struct {
	Memo struct {
		MemoType VariableLength
		MemoData VariableLength
	}
}

type Memos []Memo
