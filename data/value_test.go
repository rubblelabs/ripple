package data

import (
	. "github.com/donovanhide/ripple/testing"
	. "launchpad.net/gocheck"
)

type ValueSuite struct{}

var _ = Suite(&ValueSuite{})

var valueTests = TestSlice{
	{valueCheck("0").String(), Equals, "0", "Parse 0"},
	{valueCheck("1").String(), Equals, "1", "Parse 1"},
	{valueCheck("0.01").String(), Equals, "0.01", "Parse 0.01"},
	{valueCheck("-0").String(), Equals, "0", "Parse -0"},
	{valueCheck("-1").String(), Equals, "-1", "Parse -1"},
	{valueCheck("-0.01").String(), Equals, "-0.01", "Parse -0.01"},
	//{valueCheck("9999999999999999e80").String(), Equals, "9999999999999999e80", "Parse 9999999999999999e80"},
	{valueCheck("1e-96").String(), Equals, "0", "Parse 1e-96"},
	{ErrorCheck(NewValue("10000000000000000e80", false)), ErrorMatches, "Value overflow: .*", "Parse 10000000000000000e80"},
	{ErrorCheck(NewValue("9999999999999999e81", false)), ErrorMatches, "Value overflow: .*", "Parse 9999999999999999e81"},
	{ErrorCheck(NewValue("foo", false)), ErrorMatches, "Invalid Number: .*", "Parse foo"},

	{valueCheck("123").ZeroClone().String(), Equals, "0", "ZeroClone"},
	{valueCheck("0").IsZero(), Equals, true, "IsZero true"},
	{valueCheck("123").IsZero(), Equals, false, "IsZero false"},

	{valueCheck("-0.01").Abs().String(), Equals, "0.01", "Abs -0.01"},
	{valueCheck("0.01").Abs().String(), Equals, "0.01", "Abs 0.01"},

	{valueCheck("123").Negate().String(), Equals, "-123", "Negate 123"},
	{valueCheck("-123").Negate().String(), Equals, "123", "Negate -123"},
	{valueCheck("0").Negate().String(), Equals, "0", "Negate 0"},

	{equalValCheck("0", "0"), Equals, true, "0==0"},
	{equalValCheck("1", "0.1"), Equals, false, "1==0.1"},
	{equalValCheck("10", "0.1"), Equals, false, "10==0.1"},
	{equalValCheck("-1", "1"), Equals, false, "-1==1"},

	{addValCheck("0", "0").String(), Equals, "0", "0+0"},
	{addValCheck("0", "1").String(), Equals, "1", "0+1"},
	{addValCheck("0", "0.0001").String(), Equals, "0.0001", "0+0.0001"},
	{addValCheck("1", "0").String(), Equals, "1", "1+0"},
	{addValCheck("1", "1").String(), Equals, "2", "1+1"},
	{addValCheck("-1", "1").String(), Equals, "0", "-1+1"},
	{addValCheck("-1", "-1").String(), Equals, "-2", "-1+-1"},
	{addValCheck("1", "-1").String(), Equals, "0", "1+-1"},

	{subValCheck("0", "0").String(), Equals, "0", "0-0"},
	{subValCheck("1", "1").String(), Equals, "0", "1-1"},
	{subValCheck("-1", "0").String(), Equals, "-1", "-1-0"},
	{subValCheck("1", "-1").String(), Equals, "2", "1--1"},
	{subValCheck("0", "0.0001").String(), Equals, "-0.0001", "0-0.0001"},

	{mulValCheck("0", "0").String(), Equals, "0", "0*0"},
	{mulValCheck("1", "0").String(), Equals, "0", "1*0"},
	{mulValCheck("0", "1").String(), Equals, "0", "0*1"},
	{mulValCheck("1", "1").String(), Equals, "1", "1*1"},
	{mulValCheck("1000", "0.001").String(), Equals, "1", "1000*0.001"},
	{mulValCheck("1000", "2").String(), Equals, "2000", "1000*2"},
	{mulValCheck("1000", "-2").String(), Equals, "-2000", "1000*-2"},
	{mulValCheck("-1000", "2").String(), Equals, "-2000", "1000*-2"},
	{mulValCheck("-1000", "-2").String(), Equals, "2000", "-1000*-2"},

	{ErrorCheck(valueCheck("0").Divide(*valueCheck("0"))), ErrorMatches, "Division by zero", "0/0"},
	{divValCheck("0", "1").String(), Equals, "0", "0/1"},
	{divValCheck("1", "2").String(), Equals, "0.5", "1/2"},
	{divValCheck("-1", "2").String(), Equals, "-0.5", "-1/2"},
	{divValCheck("1", "-200").String(), Equals, "-0.005", "1/-200"},

	{valueCheck("1").Compare(*valueCheck("1")), Equals, 0, "1 Compare 1"},
	{valueCheck("0").Compare(*valueCheck("1")), Equals, -1, "0 Compare 1"},
	{valueCheck("1").Compare(*valueCheck("0")), Equals, 1, "1 Compare 0"},
	{valueCheck("0").Compare(*valueCheck("0")), Equals, 0, "0 Compare 0"},
	{valueCheck("0").Compare(*valueCheck("-1")), Equals, 1, "0 Compare -1"},
	{valueCheck("-1").Compare(*valueCheck("0")), Equals, -1, "-1 Compare 0"},
	{valueCheck("-1").Compare(*valueCheck("1")), Equals, -1, "-1 Compare 1"},
	{valueCheck("1").Compare(*valueCheck("-1")), Equals, 1, "1 Compare -1"},
	{valueCheck("1").Compare(*valueCheck("0.002")), Equals, 1, "1 Compare 0.002"},
	{valueCheck("-1").Compare(*valueCheck("0.002")), Equals, -1, "-1 Compare 0.002"},
	{valueCheck("1").Compare(*valueCheck("-0.002")), Equals, 1, "1 Compare -0.002"},
	{valueCheck("-1").Compare(*valueCheck("-0.002")), Equals, -1, "-1 Compare -0.002"},
	{valueCheck("0.002").Compare(*valueCheck("1")), Equals, -1, "0.002 Compare 1"},
	{valueCheck("-0.002").Compare(*valueCheck("1")), Equals, -1, "-0.002 Compare 1"},
	{valueCheck("0.002").Compare(*valueCheck("-1")), Equals, 1, "0.002 Compare -1"},
	{valueCheck("-0.002").Compare(*valueCheck("-1")), Equals, 1, "-0.002 Compare -1"},

	{valueCheck("1").Less(*valueCheck("1")), Equals, false, "1<1"},
	{valueCheck("0").Less(*valueCheck("1")), Equals, true, "1<1"},
}

func subValCheck(a, b string) *Value {
	if sum, err := valueCheck(a).Subtract(*valueCheck(b)); err != nil {
		panic(err)
	} else {
		return sum
	}
}

func addValCheck(a, b string) *Value {
	if sum, err := valueCheck(a).Add(*valueCheck(b)); err != nil {
		panic(err)
	} else {
		return sum
	}
}

func mulValCheck(a, b string) *Value {
	if product, err := valueCheck(a).Multiply(*valueCheck(b)); err != nil {
		panic(err)
	} else {
		return product
	}
}

func divValCheck(a, b string) *Value {
	if quotient, err := valueCheck(a).Divide(*valueCheck(b)); err != nil {
		panic(err)
	} else {
		return quotient
	}
}

func valueCheck(v string) *Value {
	if a, err := NewValue(v, false); err != nil {
		panic(err)
	} else {
		return a
	}
}

func equalValCheck(a, b string) bool {
	return valueCheck(a).Equals(*valueCheck(b))
}

func (s *ValueSuite) TestValue(c *C) {
	valueTests.Test(c)
}
