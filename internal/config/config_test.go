package config

import (
	"os"
	"testing"
)

func Test_readFromFile(t *testing.T) {

	pwd, _ := os.Getwd()
	t.Logf("pwd: %s", pwd)

	cfg := ReadFile("../../config.yml")
	if cfg == nil {
		t.Error()
	}

	t.Logf("config: %v", cfg)

	if cfg.Sub.Topic != "sub-topic" {
		t.Fail()
	}
	if cfg.Pub.Topic != "pub-topic" {
		t.Fail()
	}

}
