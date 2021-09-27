package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/latonaio/data-sweeper-kube/config"
	"bitbucket.org/latonaio/data-sweeper-kube/helper"
)

func getSweepInfo(filePath string, setting *config.Setting) (string, string, int) {
	for _, st := range setting.SweepTargets {

		for _, e := range st.FileExtention {
			if strings.HasSuffix(filePath, e) {
				return st.Name, e, st.Interval
			}
		}
	}
	return "", "", 0
}

func inSweepInterval(filePath string, interval int) bool {
	f, _ := os.Stat(filePath)

	return time.Now().Add(-(time.Millisecond * time.Duration(interval))).Before(f.ModTime())
}

func isIgnore(filePath string, setting *config.Setting) bool {
	for _, st := range setting.IgnoreMicroservices {
		if strings.Index(filePath, st.Microservice) >= 0 {
			//fmt.Println("ignore check: " + st.Microservice)

			for _, e := range st.FileExtention {
				if strings.HasSuffix(filePath, e) {
					fmt.Println("ignore: " + e + ": " + filePath)
					return true
				}
			}

			for _, n := range st.FileName {
				if strings.HasSuffix(filePath, n) {
					fmt.Println("ignore: " + n + ": " + filePath)
					return true
				}
			}
		}
	}
	return false
}

func fileDelete(filePath string, setting *config.Setting) {
	// check file type
	name, _, interval := getSweepInfo(filePath, setting)
	if name == "" {
		return
	}

	// check timestamp
	if inSweepInterval(filePath, interval) {
		return
	}

	// check ignore list
	if isIgnore(filePath, setting) {
		return
	}

	// delete this file
	fmt.Println("delete file: " + filePath)
	if err := os.Remove(filePath); err != nil {
		fmt.Println("delete failed: " + filePath)
		fmt.Println(err)
	}
}

func fileSearchRecursive(dir string, setting *config.Setting) {
	//fmt.Println("check dir: " + dir)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		if file.IsDir() {
			fileSearchRecursive(filepath.Join(dir, file.Name()), setting)
		} else {
			fileDelete(filepath.Join(dir, file.Name()), setting)
		}
	}
}

func main() {
	// 時間間隔または指定時刻を取得
	var interval time.Duration
	var alarm helper.Time
	sweepStartType := os.Getenv("SWEEP_START_TYPE")
	switch sweepStartType {
	case "interval":
		interval = time.Millisecond * 1000
		if os.Getenv("SWEEP_CHECK_INTERVAL") != "" {
			m, _ := strconv.Atoi(os.Getenv("SWEEP_CHECK_INTERVAL"))
			interval = time.Millisecond * time.Duration(m)
		}
	case "alarm":
		alarm = helper.NewTime(0, 0, 0)
		if os.Getenv("SWEEP_CHECK_ALARM") != "" {
			alarmSplice := strings.Split(os.Getenv("SWEEP_CHECK_ALARM"), ":")
			hours, _ := strconv.Atoi(alarmSplice[0])
			minutes, _ := strconv.Atoi(alarmSplice[1])
			seconds, _ := strconv.Atoi(alarmSplice[2])
			alarm = helper.NewTime(hours, minutes, seconds)

		}
	default:
		fmt.Printf("invalid sweep start type: %s", sweepStartType)
		return
	}

	baseDir := "/var/lib/aion/Data"
	if os.Getenv("AION_HOME") != "" {
		baseDir = filepath.Join(os.Getenv("AION_HOME"), "Data")
	}
	configFile := "/var/lib/aion/config/data-sweeper.yml"
	if _, err := os.Stat(configFile); err != nil {
		fmt.Println("please set configfile: " + configFile)
	}
	// 構造体作成
	c := config.GetSettingInstance()
	// 設定ファイルの読み込み
	if err := c.LoadConfig(configFile); err != nil {
		fmt.Printf("load config file error: %v", err)
		return
	}

	done := make(chan bool, 1)

	switch sweepStartType {
	case "interval":
		go func() {
			t := time.NewTicker(interval)
			for {
				select {
				case <-t.C:
					fileSearchRecursive(baseDir, c)
				case <-done:
					t.Stop()
					goto L
				}
			}
		L:
		}()
	case "alarm":
		go func() {
			for {
				// 指定時刻になると何回も繰り返し起動するためスリープを設けて防ぐ
				time.Sleep(3 * time.Second)

				t := helper.ScheduleAlarm(alarm, func() {
					fmt.Println("alarm received")
				})
				select {
				case <-t:
					fileSearchRecursive(baseDir, c)
				case <-done:
					goto L
				}
			}
		L:
		}()
	}

	server := NewServer("0.0.0.0", 8080)
	go server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	done <- true
	server.Stop()
}
