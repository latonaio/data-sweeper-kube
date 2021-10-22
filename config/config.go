package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Setting struct {
	SweepTargets        []*sweepTarget        `yaml:"sweepTargets"`
	IgnoreMicroservices []*ignoreMicroservice `yaml:"ignoreMicroservices"`
	SweepSettings       []*sweepSettings      `yaml:"sweepSettings"`
}

type ignoreMicroservice struct {
	Microservice  string   `yaml:"microservice,omitempty"`
	FileName      []string `yaml:"fileName,omitempty"`
	FileExtention []string `yaml:"fileExtention,omitempty"`
}

type sweepTarget struct {
	Name          string   `yaml:"name"`
	FileExtention []string `yaml:"fileExtention"`
	Interval      int      `yaml:"interval"`
}

type sweepSettings struct {
	SweepStartType     string `yaml:"sweepStartType"`
	SweepCheckInterval int    `yaml:"sweepCheckInterval"`
	SweepCheckAlarm    string `yaml:"sweepCheckAlarm"`
}

var sharedSettingInstance *Setting

func GetSettingInstance() *Setting {
	if sharedSettingInstance == nil {
		sharedSettingInstance = &Setting{}
	}
	return sharedSettingInstance
}

func (s *Setting) LoadConfig(configPath string) error {
	f, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, s); err != nil {
		return err
	}
	return nil
}
