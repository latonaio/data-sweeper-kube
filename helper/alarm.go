package helper

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Time is struct for representing time.
type Time struct {
	Hh int // Hours.
	Mm int // Minutes.
	Ss int // Seconds.
}

func init() {
	location := os.Getenv("TZ")
	loc, err := time.LoadLocation(location)
	if err != nil {
		loc = time.FixedZone(location, 9*60*60)
	}
	time.Local = loc
}

func NewTime(hours, minutes, seconds int) Time {
	return Time{
		Hh: hours,
		Mm: minutes,
		Ss: seconds,
	}
}

// ScheduleAlarm call this function to schedule the alarm. The callback will be called after the alarm is triggered.
func ScheduleAlarm(alarmTime Time, callback func()) (endRecSignal chan string) {
	endRecSignal = make(chan string)
	go func() {
		timeSplice := strings.Split(time.Now().Format("15:04:05"), ":")
		hh, _ := strconv.Atoi(timeSplice[0])
		mm, _ := strconv.Atoi(timeSplice[1])
		ss, _ := strconv.Atoi(timeSplice[2])

		startAlarm := getDiffSeconds(Time{hh, mm, ss}, alarmTime)

		// Setting alarm.
		time.AfterFunc(time.Duration(startAlarm)*time.Second, func() {
			callback()
			endRecSignal <- "finished recording"
			close(endRecSignal)
		})
	}()
	return
}

func getDiffSeconds(fromTime, toTime Time) int {
	fromSec := getSeconds(fromTime)
	toSec := getSeconds(toTime)
	diff := toSec - fromSec

	if diff < 0 {
		return diff + 24*60*60
	} else {
		return diff
	}
}

func getSeconds(time Time) int {
	return time.Hh*60*60 + time.Mm*60 + time.Ss
}
