// 月の指定で、入力されたUTCの日時から指定したタイムゾーンで前月のはじめの日時と最後の日時を返す
// 月または週の指定で、入力されたUTCの日付か指定したタイムゾーンで先週の最初の日時と最後の日時を返す
// なお、週を指定した場合、さらなる指定で週の初めの曜日を設定可能とする
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type Output struct {
	First string `json:"first"`
	Last  string `json:"last"`
}

// 引数はtype=month/type=week utc=2019-01-01 00:00:00のように指定する
func main() {
	tp := os.Getenv("TYPE")
	ut := os.Getenv("UTC")
	tmp := os.Getenv("WEEKDAY")
	sw, _ := strconv.Atoi(tmp)
	tz := os.Getenv("TIMEZONE")
	if tp == "" {
		log.Fatal("type is required")
	}
	if ut == "" {
		log.Fatal("utc is required")
	}
	if tz == "" {
		log.Fatal("timezone is required")
	}

	location, err := time.LoadLocation(tz)
	if err != nil {
		log.Fatal(err)
	}
	localtime, err := parseDate(&ut, location)
	if err != nil {
		log.Fatal(err)
	}
	if tp == "month" {
		first, last := getPrevMonth(localtime, location)
		out := Output{
			First: first.Format("2006-01-02"),
			Last:  last.Format("2006-01-02"),
		}
		jsonOut, err := json.Marshal(out)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonOut))
	} else if tp == "week" {
		first, last := getPrevWeek(localtime, location, sw)
		out := Output{
			First: first.Format("2006-01-02"),
			Last:  last.Format("2006-01-02"),
		}
		jsonOut, err := json.Marshal(out)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(jsonOut))
	} else {
		log.Fatal("invalid argument(type)")
		os.Exit(1)
	}
}

func parseDate(date *string, location *time.Location) (time.Time, error) {
	utc, err := time.Parse("2006-01-02 15:04:05", *date)
	if err != nil {
		return time.Time{}, err
	}
	t := utc.In(location)
	return t, nil
}

// 前月の最初の日付と最後の日付を得る
func getPrevMonth(localtime time.Time, location *time.Location) (time.Time, time.Time) {
	first := time.Date(localtime.Year(), localtime.Month()-1, 1, 0, 0, 0, 0, location)
	last := time.Date(localtime.Year(), localtime.Month(), 0, 23, 59, 59, 0, location)
	return first, last
}

// 先週の最初の日付と最後の日付を得る
func getPrevWeek(localtime time.Time, location *time.Location, sw int) (time.Time, time.Time) {
	t := time.Date(localtime.Year(), localtime.Month(), localtime.Day(), 0, 0, 0, 0, location)
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	if sw == 1 {
		first := t.AddDate(0, 0, -wd-7+1)
		last := first.Add(time.Hour * 24 * 6)
		last = last.Add((time.Hour * 23) + (time.Minute * 59) + (time.Second * 59))
		return first, last
	} else {
		first := t.AddDate(0, 0, -wd-7+sw)
		last := first.Add(time.Hour * 24 * 6)
		last = last.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		return first, last
	}
}
