package o

import (
	`time`
)

type (
	// calendars 日历月数据
	calendars []time.Time
	calendar  [2]time.Time
)

// Calendars 获取日历
func Calendars(Year int, Month time.Month) (date calendars) {
	var (
		begin = time.Date(Year, Month, 1, 0, 0, 0, 0, time.Local)
		month = time.Date(Year, Month+1, 0, 0, 0, 0, 0, time.Local)
	)
	date = append(date, begin)
	var index = 1
	for index = 1; index < month.Day(); index++ {
		date = append(date, begin.Add(time.Duration(index)*24*time.Hour))
	}
	return
}

// Calendar 日历开始结束日期
func Calendar(Year int, Month time.Month) (date [2]time.Time) {
	date = [2]time.Time{time.Date(Year, Month, 1, 0, 0, 0, 0, time.Local), time.Date(Year, Month+1, 0, 0, 0, 0, 0, time.Local)}
	return
}
