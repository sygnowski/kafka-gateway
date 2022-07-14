package util

import (
	"encoding/json"
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

func TestMatchNestedMapPath(t *testing.T) {

	nested := make(map[string]interface{})
	mr1 := make(map[string]interface{})

	mr1["row_0_1"] = "nothig else metters"
	nested["row_0_0"] = mr1

	if !MatchNestedMapPath(nested, []string{"row_0_0", "row_0_1"}) {
		t.Errorf("shouldn't fail")
	}

	if MatchNestedMapPath(nested, []string{"row_0_2", "row_0_1"}) {
		t.Errorf("should fail")
	}
}

func TestGetMapPath(t *testing.T) {

	nested := make(map[string]interface{})
	mr1 := make(map[string]interface{})

	mr1["row_0_1"] = "the one"
	nested["row_0_0"] = mr1

	if item := GetMapPath(nested, []string{"row_0_0", "row_0_1"}); item != "the one" {
		t.Errorf("shouldn't fail")
	}

	if item := GetMapPath(nested, []string{"row_0_1", "row_0_1"}); item != nil {
		t.Errorf("should fail")
	}

}

func TestGetLastMapPath(t *testing.T) {

	nested := make(map[string]interface{})
	mr1 := make(map[string]interface{})
	mr2 := make(map[string]interface{})

	mr2["row_2_0"] = "second"

	mr1["row_1_0"] = "the one"
	mr1["row_1_1"] = mr2

	nested["row_0_0"] = mr1

	full, item := GetLastMapPath(nested, []string{"row_0_0", "row_1_0"})
	if item != "the one" && !full {
		t.Errorf("shouldn't fail")
	}

	full, last := GetLastMapPath(nested, []string{"row_0_0", "row_1_1", "row_2_1"})

	m, isMap := last.(map[string]interface{})

	if full && isMap && m["row_2_0"] == "second" {
		t.Errorf("should fail")
	}
}

func TestContextWithCorrelation(t *testing.T) {
	rawJson := `
{
    "id": "A",
    "add": "B",
    "context": {
         "correlation": "dfb98431-c64b-0062-effa-9f16f52e5e31"
    }
}
`

	b := []byte(rawJson)

	dat := make(map[string]interface{})

	if err := json.Unmarshal(b, &dat); err != nil {
		t.Errorf("on test error: %s", err)
	}

	if !MatchNestedMapPath(dat, []string{"context", "correlation"}) {
		t.Errorf("correlation exists")
	}

}
