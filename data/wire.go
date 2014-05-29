package data

import (
	"encoding/binary"
	"fmt"
	"io"
)

func NewVariableByteReader(r Reader) (Reader, error) {
	if length, err := readVariableLength(r); err != nil {
		return nil, err
	} else {
		return LimitedByteReader(r, int64(length)), nil
	}
}

func unmarshalSlice(s []byte, r Reader, prefix string) error {
	n, err := r.Read(s)
	if n != len(s) {
		return fmt.Errorf("%s: short read: %d expected: %d", prefix, n, len(s))
	}
	if err != nil {
		return fmt.Errorf("%s: %s", prefix, err.Error())
	}
	return nil
}

func (v *Value) Bytes() []byte {
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

func (a *Amount) Bytes() []byte {
	return append(a.Value.Bytes(), append(a.Currency.Bytes(), a.Issuer.Bytes()...)...)
}

func (v *Value) Unmarshal(r Reader) error {
	var u uint64
	if err := binary.Read(r, binary.BigEndian, &u); err != nil {
		return err
	}
	v.Native = (u >> 63) == 0
	v.Negative = (u >> 62) == 0
	if v.Native {
		v.Num = u & ((1 << 62) - 1)
		v.Offset = 0
	} else {
		v.Num = u & ((1 << 54) - 1)
		v.Offset = int64((u>>54)&((1<<8)-1)) - 97
	}
	return nil
}

func (v *Value) Marshal(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, v.Bytes())
}

func (a *Amount) Unmarshal(r Reader) error {
	a.Value = new(Value)
	if err := a.Value.Unmarshal(r); err != nil {
		return err
	}
	if a.Value.Native {
		return nil
	}
	if err := unmarshalSlice(a.Currency[:], r, "Currency"); err != nil {
		return err
	}
	if err := unmarshalSlice(a.Issuer[:], r, "Issuer"); err != nil {
		return err
	}
	return nil
}

func (a *Amount) Marshal(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, a.Bytes())
}

func (c *Currency) Unmarshal(r Reader) error {
	return unmarshalSlice(c[:], r, "Currency")
}

func (c *Currency) Marshal(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, c.Bytes())
}

func (h *Hash128) Unmarshal(r Reader) error {
	return unmarshalSlice(h[:], r, "Hash128")
}

func (h *Hash128) Marshal(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, h.Bytes())
}

func (h *Hash160) Unmarshal(r Reader) error {
	return unmarshalSlice(h[:], r, "Hash160")
}

func (h *Hash160) Marshal(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, h.Bytes())
}

func (h *Hash256) Unmarshal(r Reader) error {
	return unmarshalSlice(h[:], r, "Hash256")
}

func (h *Hash256) Marshal(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, h.Bytes())
}

func writeVariableLength(w io.Writer, b []byte) error {
	n := len(b)
	var err error
	switch {
	case n < 0 || n > 918744:
		return fmt.Errorf("Unsupported Variable Length encoding: %d", n)
	case n <= 192:
		err = binary.Write(w, binary.BigEndian, uint8(n))
	case n <= 12480:
		n -= 193
		err = binary.Write(w, binary.BigEndian, [2]uint8{193 + uint8(n>>8), uint8(n)})
	case n <= 918744:
		n -= 12481
		v := [3]uint8{uint8(241 + (n >> 16)), uint8(n >> 8), uint8(n)}
		err = binary.Write(w, binary.BigEndian, v)
	}
	if err != nil {
		return err
	}
	return binary.Write(w, binary.BigEndian, b)
}

func readVariableLength(r Reader) (int, error) {
	var first, second, third byte
	var err error
	if first, err = r.ReadByte(); err != nil {
		return 0, err
	}
	switch {
	case first <= 192:
		return int(first), nil
	case first <= 240:
		if second, err = r.ReadByte(); err != nil {
			return 0, nil
		}
		return 193 + int(first-193)*256 + int(second), nil
	case first <= 254:
		if second, err = r.ReadByte(); err != nil {
			return 0, nil
		}
		if third, err = r.ReadByte(); err != nil {
			return 0, nil
		}
		return 12481 + int(first-241)*65536 + int(second)*256 + int(third), nil
	}
	return 0, fmt.Errorf("Unsupported Variable Length encoding")
}

func (v *Vector256) Unmarshal(r Reader) error {
	length, err := readVariableLength(r)
	if err != nil {
		return err
	}
	count := length / 32
	*v = make(Vector256, count)
	for i := 0; i < count; i++ {
		if err := (*v)[i].Unmarshal(r); err != nil {
			return err
		}
	}
	return nil
}

func (v *Vector256) Marshal(w io.Writer) error {
	var b []byte
	for _, h := range *v {
		b = append(b, h[:]...)
	}
	return writeVariableLength(w, b)
}

func (v *VariableLength) Unmarshal(r Reader) error {
	length, err := readVariableLength(r)
	if err != nil {
		return err
	}
	*v = make(VariableLength, length)
	return unmarshalSlice(*v, r, "VariableLength")
}

func (v *VariableLength) Marshal(w io.Writer) error {
	return writeVariableLength(w, v.Bytes())
}

func readExpectedLength(r Reader, dest []byte, prefix string) error {
	length, err := readVariableLength(r)
	switch {
	case err != nil:
		return fmt.Errorf("%s: %s", prefix, err.Error())
	case length == 0:
		return nil
	case length == len(dest):
		return unmarshalSlice(dest, r, prefix)
	default:
		return fmt.Errorf("%s: wrong length %d expected: %d", prefix, length, len(dest))
	}
}

func (a *Account) Unmarshal(r Reader) error {
	return readExpectedLength(r, a[:], "Account")
}

func (a *Account) Marshal(w io.Writer) error {
	return writeVariableLength(w, a.Bytes())
}

func (k *PublicKey) Unmarshal(r Reader) error {
	return readExpectedLength(r, k[:], "PublicKey")
}

func (k *PublicKey) Marshal(w io.Writer) error {
	return writeVariableLength(w, k.Bytes())
}

func (k *RegularKey) Unmarshal(r Reader) error {
	return readExpectedLength(r, k[:], "RegularKey")
}

func (k *RegularKey) Marshal(w io.Writer) error {
	return writeVariableLength(w, k.Bytes())
}

func (p *Paths) Unmarshal(r Reader) error {
	for i := 0; ; i++ {
		*p = append(*p, []Path{})
		for entry, err := r.ReadByte(); entry != 0xFF; entry, err = r.ReadByte() {
			if err != nil {
				return err
			}
			if entry == 0x00 {
				return nil
			}
			var path Path
			if entry&0x01 > 0 {
				path.Account = new(Account)
				if _, err := r.Read(path.Account[:]); err != nil {
					return err
				}
			}
			if entry&0x10 > 0 {
				path.Currency = new(Currency)
				if _, err := r.Read(path.Currency[:]); err != nil {
					return err
				}
			}
			if entry&0x20 > 0 {
				path.Issuer = new(Account)
				if _, err := r.Read(path.Issuer[:]); err != nil {
					return err
				}
			}
			(*p)[i] = append((*p)[i], path)
		}
	}
}

func (p *Paths) Marshal(w io.Writer) error {
	return nil
}

func (m *Memos) Unmarshal(r Reader) error {
	return nil
}

func (m *Memos) Marshal(w io.Writer) error {
	return nil
}

func (e NodeEffects) Unmarshal(r Reader) error {
	return nil
}

func (e NodeEffects) Marshal(w io.Writer) error {
	return nil
}
