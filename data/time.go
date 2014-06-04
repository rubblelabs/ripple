package data

import (
	"time"
)

const (
	rippleTimeEpoch  int64  = 946684800
	rippleTimeFormat string = "2006-Jan-02 15:04:05"
)

type RippleTime struct {
	T uint32
}

type RippleHumanTime struct {
	*RippleTime
}

func NewRippleTime(t uint32) *RippleTime {
	return &RippleTime{t}
}

func convertToRippleTime(t time.Time) uint32 {
	return uint32(t.Sub(time.Unix(rippleTimeEpoch, 0)).Nanoseconds() / 1000000000)
}

func (t *RippleTime) time() time.Time {
	return time.Unix(int64(t.T)+rippleTimeEpoch, 0)
}

func Now() *RippleTime {
	return &RippleTime{convertToRippleTime(time.Now())}
}

func (t *RippleTime) SetString(s string) error {
	v, err := time.Parse(rippleTimeFormat, s)
	if err != nil {
		return err
	}
	t.SetUint32(convertToRippleTime(v))
	return nil
}

func (t *RippleTime) SetUint32(n uint32) {
	t.T = n
}

func (t *RippleTime) Uint32() uint32 {
	return t.T
}

func (t *RippleTime) Human() *RippleHumanTime {
	return &RippleHumanTime{t}
}

func (t *RippleTime) String() string {
	return t.time().UTC().Format(rippleTimeFormat)
}

func (t *RippleTime) Short() string {
	return t.time().UTC().Format("15:04:05")
}
