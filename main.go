package main

import (
	"os"
	"time"
	"fmt"
	"strings"
	"strconv"
	"io/ioutil"
	"path/filepath"

	"bitbucket.org/latonaio/data-sweeper-kube/config"
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

    return time.Now().Add(-(time.Millisecond*time.Duration(interval))).Before(f.ModTime())
}

func isIgnore(filePath string, setting *config.Setting) bool {
	for _, st := range setting.IgnoreMicroservices {
		if strings.Index(filePath, st.Microservice) >= 0 {
			fmt.Println("ignore check: "+st.Microservice)

			for _, e := range st.FileExtention {
				if strings.HasSuffix(filePath, e) {
					fmt.Println("ignore: "+e+": "+filePath)
					return true
				}
			}

			for _, n := range st.FileName {
				if strings.HasSuffix(filePath, n) {
					fmt.Println("ignore: "+n+": "+filePath)
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
	if isIgnore(filePath, setting){
		return
	}

	// delete this file
	fmt.Println("delete file: "+filePath)
	if err := os.Remove(filePath); err != nil {
		fmt.Println("delete failed: "+filePath)
		fmt.Println(err)
	}
}

func fileSearchRecursive(dir string, setting *config.Setting) {
	fmt.Println("check dir: "+dir)
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
	interval := time.Millisecond * 1000
	if os.Getenv("SWEEP_CHECK_INTERVAL") != "" {
        m, _ := strconv.Atoi(os.Getenv("SWEEP_CHECK_INTERVAL"))
		interval = time.Millisecond * time.Duration(m)
	}
	baseDir := "/var/lib/aion/Data"
	if os.Getenv("AION_HOME") != "" {
		baseDir = filepath.Join(os.Getenv("AION_HOME"), "Data")
	}
	configFile := "/var/lib/aion/config/data-sweeper.yml"

	if _, err := os.Stat(configFile); err != nil {
		fmt.Println("please set configfile: "+configFile)
	}

    c := config.GetSettingInstance() // 構造体の取得
    s, _ := c.LoadConfig(configFile) // 設定ファイルの読み込み
	t := time.NewTicker(interval)
	for {
		select {
		case <-t.C:
			fileSearchRecursive(baseDir, s)
		}
	}
	t.Stop()

}
