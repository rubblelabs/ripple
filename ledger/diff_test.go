package ledger

import (
	"bytes"
	"github.com/donovanhide/ripple/data"
	"github.com/donovanhide/ripple/storage"
	"testing"
)

func checkErr(t *testing.T, err error) {
	if err != nil {
		t.Fatalf(err.Error())
	}
}

var expectedDiff = `0,A,0,2C23D15B6B549123,InnerNode,525
0,D,0,AF47E9E91A41621B,InnerNode,525
0,A,1,29FD2F34869B2E46,InnerNode,525
0,D,1,271E1B9B1B1FB8C7,InnerNode,525
0,A,1,067A065323B98104,InnerNode,525
0,D,1,724CA5CAEB55D794,InnerNode,525
0,A,2,C97303390D8FF71B,AccountRoot,132
0,D,2,2A75953DB729CC20,AccountRoot,132
0,A,2,FFB92B95013668AE,InnerNode,525
0,D,2,C62002973CAB176F,InnerNode,525
0,A,3,0E386F1549BD2B10,LedgerHashes,8261
0,D,3,28372696D26A27D4,LedgerHashes,8261
`

func TestDiff(t *testing.T) {
	mem, err := storage.NewMemoryDB("testdata/38129-32570.gz")
	checkErr(t, err)
	first, err := data.NewHash256("2C23D15B6B549123FB351E4B5CDE81C564318EB845449CD43C3EA7953C4DB452")
	checkErr(t, err)
	second, err := data.NewHash256("AF47E9E91A41621B0F8AC5A119A5AD8B9E892147381BEAF6F2186127B89A44FF")
	checkErr(t, err)
	diff, err := Diff(*first, *second, mem)
	checkErr(t, err)
	var buf bytes.Buffer
	checkErr(t, diff.Dump(uint32(0), &buf))
	//FIXME:
	// if buf.String() != expectedDiff {
	// 	t.Errorf("Wrong diff!\n%s\nExpected:\n%s\n", buf.String(), expectedDiff)
	// }
}
