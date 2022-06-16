// You can edit this code!
// Click here and start typing.
package main

import (
	"fmt"
	"time"
)

func main() {
	arr := []interface{}{1.3, 1, 2}
	for _, v := range arr {
		switch v.(type) {
		case int:
			fmt.Println("==> int", v.(int))
		case float64:
			fmt.Println("==> float64", v.(float64))
		}
	}

	fmt.Println(fmt.Sprintf(`%s %% 86400000 `, "aaa"))

	var a interface{}
	a = 1000.1
	switch a.(type) {
	case int:
	}
	b, ok := a.(float64)
	fmt.Println(b, ok, fmt.Sprintf("%s", a))

	fmt.Println("Hello, 世界")

	t1, _ := time.ParseInLocation("2006-01-02", "2022-04-06", time.FixedZone("GST", 8*3600))
	fmt.Println(t1, t1.UnixMilli())

	t2, _ := time.ParseInLocation("2006-01-02", "2022-04-06", time.FixedZone("GST", 6*3600))
	fmt.Println(t2, t2.UnixMilli())

	fmt.Println(t2.UnixMilli() - t1.UnixMilli())

	now := time.Now().In(time.FixedZone("GST", 6*3600))
	z, o := now.Zone()
	fmt.Println("==> now.Location", now.Location(), z, o)
	fmt.Println(time.Now().UnixMilli())
	fmt.Println(now.UnixMilli())

	fmt.Println(time.Now())
	fmt.Println(time.Now().In(time.FixedZone("GST", 6*3600)))

	d := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.FixedZone("GST", 6*3600))
	fmt.Println(d)

	d2 := d.AddDate(0, 0, 7)
	fmt.Println(d2)

	t := time.Now()
	zone, offset := t.Zone()
	fmt.Println(zone, offset)

	loc, _ := time.LoadLocation("America/New_York")
	t2 = time.Now().In(loc)
	fmt.Println(t2.Zone())

	aa := []int{}
	fmt.Println("===> ", aa[1:])
}
