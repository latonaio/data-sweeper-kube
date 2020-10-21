package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type Setting struct {
	SweepTargets        []*SweepTarget        `yaml:"sweepTargets"`
	IgnoreMicroservices []*ignoreMicroservice `yaml:"ignoreMicroservices"`
}

type ignoreMicroservice struct {
	Microservice  string   `yaml:"microservice,omitempty"`
	FileName      []string `yaml:"fileName,omitempty"`
	FileExtention []string `yaml:"fileExtention,omitempty"`
}

type SweepTarget struct {
	Name          string   `yaml:"name"`
	FileExtention []string `yaml:"fileExtention"`
	Interval      int      `yaml:"interval"`
}

var sharedSettingInstance *Setting

func GetSettingInstance() *Setting {
	if sharedSettingInstance == nil {
		sharedSettingInstance = &Setting{}
	}
	return sharedSettingInstance
}

func (s *Setting) LoadConfig(configPath string) (*Setting, error) {
	f, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, s); err != nil {
		return nil, err
	}
	return s, nil
}
