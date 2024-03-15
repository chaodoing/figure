package o

import (
	"time"
)

// Unix 时间戳转换
type Unix int64

func (e Unix) MarshalText() ([]byte, error) {
	if e != 0 {
		tm := time.Unix(int64(e), 0)
		t := tm.Format(FORMAT_DATE)
		return []byte(t), nil
	} else {
		return []byte("0"), nil
	}
}

func (e *Unix) UnmarshalText(data []byte) error {
	var (
		ss  = string(data)
		err error
		ts  time.Time
	)
	switch len(ss) {
	case len(FORMAT_DATE_TIME):
		ts, err = time.Parse(FORMAT_DATE_TIME, string(data))
	case len(FORMAT_MONTH):
		ts, err = time.Parse(FORMAT_MONTH, string(data))
	case len(FORMAT_DATE):
		ts, err = time.Parse(FORMAT_DATE, string(data))
	case len(FORMAT_TIME):
		ts, err = time.Parse(FORMAT_TIME, string(data))
	}
	if err == nil {
		*e = Unix(ts.Unix())
	}
	return err
}
