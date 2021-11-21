package timex

import "time"

func GetNowUTCMilli() int64 {
	return time.Now().UTC().UnixNano() / 1e6
}
