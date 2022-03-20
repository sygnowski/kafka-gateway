package util

import (
	"regexp"
	"testing"
)

func TestUUID(t *testing.T) {
	uuid := UUID()

	match, _ := regexp.MatchString("\\b[0-9a-f]{8}\\b-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-\\b[0-9a-f]{12}\\b", uuid)
	if !match {
		t.Errorf("UUID (%s) not valid\n", uuid)
	}

}
