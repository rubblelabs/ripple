package data

import (
	"fmt"
	"strconv"
	"time"
)

const (
	rippleEpoch int64 = 946684800
)

type RippleTime time.Time

func NewRippleTime(t uint32) *RippleTime {
	r := time.Unix(int64(t)+rippleEpoch, 0)
	return (*RippleTime)(&r)
}

func Now() int64 {
	return time.Now().Sub(time.Unix(rippleEpoch, 0)).Nanoseconds() / 1000000000
}

func (t *RippleTime) UnmarshalJSON(buf []byte) error {
	if unix, err := strconv.ParseInt(string(buf), 10, 64); err != nil {
		return fmt.Errorf("Bad RippleTime:%s", string(buf))
	} else {
		*t = RippleTime(time.Unix(unix+rippleEpoch, 0))
	}
	return nil
}

func (t *RippleTime) String() string {
	return time.Time(*t).Format("2006-01-02 15:04:05")
}

func (t *RippleTime) Short() string {
	return time.Time(*t).Format("15:04:05")
}
