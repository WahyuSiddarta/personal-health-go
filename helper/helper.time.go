package helper

import "time"

// TimeIn :
func TimeInWIB(t time.Time) (time.Time, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err == nil {
		t = t.In(loc)
	}
	return t, err
}
