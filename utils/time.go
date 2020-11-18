package utils

import "time"

var (
	jst *time.Location
)

func init() {
	jst = time.FixedZone("Asia/Tokyo", 9*60*60)
}

//NowJST -
func NowJST() time.Time {
	return time.Now().UTC().In(jst)
}

//FormatTime -
func FormatTime(t time.Time, f string) string {
	if f == "" {
		f = "2006/01/02 15:04:05"
	}

	return t.In(jst).Format(f)
}
