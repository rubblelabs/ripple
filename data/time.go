package data

import (
	"time"
)

const (
	rippleTimeEpoch  int64  = 946684800
	rippleTimeFormat string = "2006-Jan-02 15:04:05"
)

type RippleTime time.Time

func NewRippleTime(t uint32) *RippleTime {
	r := time.Unix(int64(t)+rippleTimeEpoch, 0)
	return (*RippleTime)(&r)
}

func (t *RippleTime) Parse(s string) error {
	p, err := time.Parse(rippleTimeFormat, s)
	if err != nil {
		return err
	}
	*t = RippleTime(p)
	return nil
}

func Now() int64 {
	return time.Now().Sub(time.Unix(rippleTimeEpoch, 0)).Nanoseconds() / 1000000000
}

func (t *RippleTime) String() string {
	return time.Time(*t).UTC().Format(rippleTimeFormat)
}

func (t *RippleTime) Short() string {
	return time.Time(*t).UTC().Format("15:04:05")
}
