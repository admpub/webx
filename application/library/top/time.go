package top

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"time"
)

var (
	ErrInvalidDuration = errors.New("invalid duration")
	DurationDay        = time.Hour * 24
	DurationMonth      = DurationDay * 30
	DurationWeek       = DurationDay * 7
	DurationYear       = DurationDay * 365
	MonthMaxSeconds    = 31 * 86400 // 月份中可能的最大秒数
)
var durationRegex = regexp.MustCompile(`^[\d]+([.][\d]+)?(ns|us|µs|ms|s|m|h|d|w|mo|y)$`)

// ParseDuration 解析持续时间(在支持标准库time.ParseDuration的基础上增加了年(y)月(mo)周(w)日(d)的支持)
func ParseDuration(s string) (time.Duration, error) {
	if !durationRegex.MatchString(s) {
		return 0, ErrInvalidDuration
	}
	size := len(s)
	if size > 2 && s[size-2:] == `mo` {
		n, err := strconv.ParseInt(s[0:size-2], 10, 64)
		if err != nil {
			return 0, err
		}
		return time.Until(time.Now().AddDate(0, int(n), 0)), nil
		//return DurationMonth * time.Duration(n), nil
	}
	switch s[size-1] {
	case 'd':
		n, err := strconv.ParseInt(s[0:size-1], 10, 64)
		if err != nil {
			return 0, err
		}
		return DurationDay * time.Duration(n), nil
	case 'w':
		n, err := strconv.ParseInt(s[0:size-1], 10, 64)
		if err != nil {
			return 0, err
		}
		return DurationWeek * time.Duration(n), nil
	case 'y':
		n, err := strconv.ParseInt(s[0:size-1], 10, 64)
		if err != nil {
			return 0, err
		}
		return DurationYear * time.Duration(n), nil
	}
	return time.ParseDuration(s)
}

// IsSameDay 是否为同一天
func IsSameDay(last time.Time, nows ...time.Time) bool {
	var now time.Time
	if len(nows) > 0 {
		now = nows[0]
	}
	if now.IsZero() {
		now = time.Now()
	}
	return last.Day() == now.Day() && math.Abs(float64(now.Unix()-last.Unix())) < 86400
}

// MonthDay 计算某个月的天数
func MonthDay(year int, month int) (days int) {
	if month != 2 {
		if month == 4 || month == 6 || month == 9 || month == 11 {
			days = 30
		} else {
			days = 31
		}
		return
	}
	if ((year%4) == 0 && (year%100) != 0) || (year%400) == 0 {
		days = 29
	} else {
		days = 28
	}
	return
}

// MonthDayByTime 计算某个月的天数
func MonthDayByTime(t time.Time) int {
	return MonthDay(t.Year(), int(t.Month()))
}

// IsSameMonth 是否为同一月
func IsSameMonth(last time.Time, nows ...time.Time) bool {
	var now time.Time
	if len(nows) > 0 {
		now = nows[0]
	}
	if now.IsZero() {
		now = time.Now()
	}
	return last.Month() == now.Month() && math.Abs(float64(now.Unix()-last.Unix())) < float64(MonthMaxSeconds)
}

func TodayTimestamp() (startTs int64, endTs int64) {
	now := time.Now()
	startTs = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).Unix()
	endTs = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local).Unix()
	return
}
