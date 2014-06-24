package crypto

import (
	"github.com/donovanhide/ripple/testing"
	. "launchpad.net/gocheck"
)

type PowSuite struct{}

var _ = Suite(&PowSuite{})

var powTests = []struct {
	challenge  string
	target     string
	solution   string
	iterations uint32
	slow       bool
}{
	{
		// Low difficulty
		"0BC4775DC5A00A5EF80022C1562D6CED655942886F7CE12D76620164C6A613B5",
		"0CFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"385CCC3C47E68AAFA5FE8D57883FE060E43F90B734217C39E81930EF67216623",
		98304,
		false,
	},
	{
		// Difficulty 5
		"F5FC1C5BA0AEAB99AE545F4D7B2256807C1E4014F044D6BAA9BD5A04838F8850",
		"0CFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"326EAF28DF40C3D7B23793905C236EBBEB03653E0767B569CBEDD84763F9DBF5",
		262144,
		true,
	},
	{
		// Diffculty 8
		"7BE5ABA26C95A0E21FF6E25C6DD298B72AB7D3559FD26578A83A5BA2E77AEE77",
		"07FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
		"D777413E7D85060E3139B1C429B1B97B16378BC42736A7C60B94CEE2876DCF5D",
		393216,
		true,
	},
}

func (p *PowSuite) TestProofOfWork(c *C) {
	for _, t := range powTests {
		if !*testing.RunSlow && t.slow {
			continue
		}
		pow := NewProofOfWork(hexToBytes(t.challenge), hexToBytes(t.target), t.iterations)
		nonce := pow.Solve()
		found := pow.Check(nonce)
		c.Check(found, Equals, true)
		solution := pow.Check(hexToBytes(t.solution))
		c.Check(solution, Equals, true)
	}
}
