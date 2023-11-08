package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Cfg struct {
	App CfgApp  `yaml:"app"`
	Pub CfgSpec `yaml:"pub"`
	Sub CfgSpec `yaml:"sub"`
}

type CfgApp struct {
	Port    int
	Timeout uint32
	Context string
}

type CfgSpec struct {
	Topic      string
	Properties []string //`yaml:",flow"`
}

func ReadConfig() *Cfg {

	confTxt := os.Getenv("APP_CONFIG")

	switch {
	case confTxt != "":
		return parse([]byte(confTxt))
	default:
		log.Panicf("Missing configuration setup")
		return nil
	}

}

func ReadFile(path string) *Cfg {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Panicf("Missing config file")
		return nil
	}
	return parse([]byte(content))

}

func parse(data []byte) *Cfg {
	config := Cfg{}
	err := yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Invalid Config: %v", err)
	}

	return &config
}
