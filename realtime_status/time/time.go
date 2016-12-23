package time

import (
	"time"
)

func GetNowTime() (int64, string) {
	c_time := time.Now().Unix()
	var timeLayout = "2006-01-02 15:04:05"
	var dataTimeStr = time.Unix(c_time, 0).Format(timeLayout)
	return c_time, dataTimeStr
}
