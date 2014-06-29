package data

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
)

const (
	minOffset        int64  = -96
	maxOffset        int64  = 80
	minValue         uint64 = 1000000000000000
	maxValue         uint64 = 9999999999999999
	maxNative        uint64 = 9000000000000000000
	maxNativeNetwork uint64 = 100000000000000000
	notNative        uint64 = 0x8000000000000000
	positive         uint64 = 0x4000000000000000
	maxNativeSqrt    uint64 = 3000000000
	maxNativeDiv     uint64 = 2095475792 // MaxNative / 2^32
	tenTo14          uint64 = 100000000000000
	tenTo14m1        uint64 = tenTo14 - 1
	tenTo17          uint64 = tenTo14 * 1000
	tenTo17m1        uint64 = tenTo17 - 1
	xrpPrecision     uint64 = 1000000
)

var (
	bigTen        = big.NewInt(10)
	bigTenTo14    = big.NewInt(0).SetUint64(tenTo14)
	bigTenTo17    = big.NewInt(0).SetUint64(tenTo17)
	zeroNative    Value
	zeroNonNative Value
	xrpMultipler  = newValue(true, false, xrpPrecision, 0)
)

type Value struct {
	Native   bool
	Negative bool
	Num      uint64
	Offset   int64
}

func init() {
	zeroNative.Native = true
	zeroNonNative.Native = false
	if err := zeroNative.canonicalise(); err != nil {
		panic(err)
	}
	if err := zeroNonNative.canonicalise(); err != nil {
		panic(err)
	}
	if err := xrpMultipler.canonicalise(); err != nil {
		panic(err)
	}
}

func newValue(native, negative bool, num uint64, offset int64) *Value {
	return &Value{
		Native:   native,
		Negative: negative,
		Num:      num,
		Offset:   offset,
	}
}

// NewNativeValue returns a Value with n drops.
// If the value is impossible an error is returned.
func NewNativeValue(n int64) (*Value, error) {
	v := newValue(true, n < 0, uint64(n), 0)
	return v, v.canonicalise()
}

// Match fields:
// 0 = whole input
// 1 = sign
// 2 = integer portion
// 3 = whole fraction (with '.')
// 4 = fraction (without '.')
// 5 = whole exponent (with 'e')
// 6 = exponent sign
// 7 = exponent number
var valueRegex = regexp.MustCompile("([+-]?)(\\d*)(\\.(\\d*))?([eE]([+-]?)(\\d+))?")

func NewValue(s string, native bool) (*Value, error) {
	var err error
	v := Value{
		Native: native,
	}
	matches := valueRegex.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("Invalid Number: %s", s)
	}
	if len(matches[2])+len(matches[4]) > 32 {
		return nil, fmt.Errorf("Overlong Number: %s", s)
	}
	if matches[1] == "-" {
		v.Negative = true
	}
	if len(matches[4]) == 0 {
		if v.Num, err = strconv.ParseUint(matches[2], 10, 64); err != nil {
			return nil, fmt.Errorf("Invalid Number: %s", s)
		}
		v.Offset = 0
	} else {
		if v.Num, err = strconv.ParseUint(matches[2]+matches[4], 10, 64); err != nil {
			return nil, fmt.Errorf("Invalid Number: %s", s)
		}
		v.Offset = -int64(len(matches[4]))
	}
	if len(matches[5]) > 0 {
		exp, err := strconv.ParseInt(matches[7], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("Invalid Number: %s", s)
		}
		if matches[6] == "-" {
			v.Offset -= exp
		} else {
			v.Offset += exp
		}
	}
	if v.Native && len(matches[3]) > 0 {
		v.Offset += 6
	}
	return &v, v.canonicalise()
}

func (v *Value) canonicalise() error {
	if v.Native {
		if v.Num == 0 {
			v.Offset = 0
			v.Negative = false
		} else {
			for v.Offset < 0 {
				v.Num /= 10
				v.Offset++
			}
			for v.Offset > 0 {
				v.Num *= 10
				v.Offset--
			}
			if v.Num > maxNative {
				return fmt.Errorf("Native amount out of range: %s", v.debug())
			}
		}
	} else {
		if v.Num == 0 {
			v.Offset = -100
			v.Negative = false
		} else {
			for v.Num < minValue && v.Offset > minOffset {
				v.Num *= 10
				v.Offset--
			}
			for v.Num > maxValue {
				if v.Offset >= maxOffset {
					return fmt.Errorf("Value overflow: %s", v.debug())
				}
				v.Num /= 10
				v.Offset++
			}
			if v.Offset < minOffset || v.Num < minValue {
				v.Num = 0
				v.Offset = 0
				v.Negative = false
			}
			if v.Offset > maxOffset {
				return fmt.Errorf("Value overflow: %s", v.debug())
			}
		}
	}
	return nil
}

// Clone returns a Value which is a copy of v.
func (v Value) Clone() *Value {
	return newValue(v.Native, v.Negative, v.Num, v.Offset)
}

// ZeroClone returns a zero Value, native or non-native depending on v's setting.
func (v Value) ZeroClone() *Value {
	if v.Native {
		return &zeroNative
	}
	return &zeroNonNative
}

// Abs returns a copy of v with a positive sign
func (v Value) Abs() *Value {
	return newValue(v.Native, false, v.Num, v.Offset)
}

// Negate returns a new Value with the opposite sign of v.
func (v Value) Negate() *Value {
	return newValue(v.Native, !v.Negative, v.Num, v.Offset)
}

func (a Value) factor(b Value) (int64, int64, int64) {
	ao, bo := a.Offset, b.Offset
	av, bv := int64(a.Num), int64(b.Num)
	if a.Negative {
		av = -av
	}
	if b.Negative {
		bv = -bv
	}
	for ; ao < bo; ao++ {
		av /= 10
	}
	for ; bo < ao; bo++ {
		bv /= 10
	}
	return av, bv, ao
}

// Add adds a to b and returns the sum as a new Value.
func (a Value) Add(b Value) (*Value, error) {
	switch {
	case a.Native != b.Native:
		return nil, fmt.Errorf("Cannot add native and non-native values")
	case a.IsZero():
		return b.Clone(), nil
	case b.IsZero():
		return a.Clone(), nil
	default:
		av, bv, ao := a.factor(b)
		v := newValue(a.Native, (av+bv) < 0, abs(av+bv), ao)
		return v, v.canonicalise()
	}
}

func (a Value) Subtract(b Value) (*Value, error) {
	return a.Add(*b.Negate())
}

func normalise(a, b Value) (uint64, uint64, int64, int64) {
	av, bv := a.Num, b.Num
	ao, bo := a.Offset, b.Offset
	if a.Native {
		for ; av < minValue; av *= 10 {
			ao--
		}
	}
	if b.Native {
		for ; bv < minValue; bv *= 10 {
			bo--
		}
	}
	return av, bv, ao, bo
}

func (a Value) Multiply(b Value) (*Value, error) {
	if a.IsZero() || b.IsZero() {
		return a.ZeroClone(), nil
	}
	if a.Native && b.Native {
		min := min64(a.Num, b.Num)
		max := max64(a.Num, b.Num)
		if min > maxNativeSqrt || (((max >> 32) * min) > maxNativeDiv) {
			return nil, fmt.Errorf("Native value overflow: %s*%s", a.debug(), b.debug())
		}
		return NewNativeValue(int64(min * max))
	}
	av, bv, ao, bo := normalise(a, b)
	// Compute (numerator * denominator) / 10^14 with rounding
	// 10^16 <= result <= 10^18
	m := big.NewInt(0).SetUint64(av)
	m.Mul(m, big.NewInt(0).SetUint64(bv))
	m.Div(m, bigTenTo14)
	// 10^16 <= product <= 10^18
	if len(m.Bytes()) > 64 {
		return nil, fmt.Errorf("Multiply: %s*%s", a.debug(), b.debug())
	}
	v := newValue(a.Native, a.Negative != b.Negative, m.Uint64()+7, ao+bo+14)
	return v, v.canonicalise()
}

func (num Value) Divide(den Value) (*Value, error) {
	if den.IsZero() {
		return nil, fmt.Errorf("Division by zero")
	}
	if num.IsZero() {
		return num.ZeroClone(), nil
	}
	av, bv, ao, bo := normalise(num, den)
	// Compute (numerator * 10^17) / denominator
	d := big.NewInt(0).SetUint64(av)
	d.Mul(d, bigTenTo17)
	d.Div(d, big.NewInt(0).SetUint64(bv))
	// 10^16 <= quotient <= 10^18
	if len(d.Bytes()) > 64 {
		return nil, fmt.Errorf("Divide: %s/%s", num.debug(), den.debug())
	}
	v := newValue(num.Native, num.Negative != den.Negative, d.Uint64()+5, ao-bo-17)
	return v, v.canonicalise()
}

func (num Value) Ratio(den Value) (*Value, error) {
	quotient, err := num.Divide(den)
	if err != nil {
		return nil, err
	}
	if den.Native {
		return quotient.Multiply(*xrpMultipler)
	}
	return quotient, nil
}

// Less compares values and returns true
// if v is less than other
func (a Value) Less(b Value) bool {
	return a.Compare(b) < 0
}

func (a Value) Equals(b Value) bool {
	return a.Native == b.Native && a.Compare(b) == 0
}

//Compare returns an integer comparing two Values. The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
func (a Value) Compare(b Value) int {
	switch {
	case a.Negative != b.Negative, a.Offset > b.Offset:
		if a.Negative {
			return -1
		}
		return 1
	case a.Offset < b.Offset:
		if a.Negative {
			return 1
		}
		return -1
	}
	switch {
	case a.Num > b.Num:
		if a.Negative {
			return -1
		}
		return 1
	case a.Num < b.Num:
		if a.Negative {
			return 1
		}
		return -1
	default:
		return 0
	}
}

// Indicates when value should be String()ed in scientific notation.
func (v Value) isScientific() bool {
	return v.Offset != 0 && (v.Offset < -25 || v.Offset > -5)
}

func (v Value) IsZero() bool {
	return v.Num == 0
}

func (v Value) Bytes() []byte {
	var u uint64
	if !v.Negative {
		u |= 1 << 62
	}
	if !v.Native {
		u |= 1 << 63
		u |= v.Num & ((1 << 54) - 1)
		u |= uint64(v.Offset+97) << 54
	} else {
		u |= v.Num & ((1 << 62) - 1)
	}
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], u)
	return b[:]
}

func (v Value) String() string {
	if v.IsZero() {
		return "0"
	}
	if !v.Native && v.isScientific() {
		value := strconv.FormatUint(v.Num, 10)
		origLen := len(value)
		value = strings.TrimRight(value, "0")
		offset := strconv.FormatInt(v.Offset+int64(origLen-len(value)), 10)
		if v.Negative {
			return "-" + value + "e" + offset
		}
		return value + "e" + offset
	}
	value := big.NewInt(int64(v.Num))
	if v.Negative {
		value.Neg(value)
	}
	var offset *big.Int
	if v.Native {
		offset = big.NewInt(-6)
	} else {
		offset = big.NewInt(v.Offset)
	}
	exp := offset.Exp(bigTen, offset.Neg(offset), nil)
	rat := big.NewRat(0, 1).SetFrac(value, exp)
	left := rat.FloatString(0)
	if rat.IsInt() {
		return left
	}
	length := len(left)
	if v.Negative {
		length -= 1
	}
	return strings.TrimRight(rat.FloatString(32-length), "0")
}

func (v Value) debug() string {
	return fmt.Sprintf("Native: %t Negative: %t Value: %d Offset: %d", v.Native, v.Negative, v.Num, v.Offset)
}
