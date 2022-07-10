package kafint

import (
	"encoding/json"
	"testing"

	"s7i.io/kafka-gateway/internal/util"
)

func Test_attachContext(t *testing.T) {
	input := `{"id": "123"}`

	ki := KafkaIntegrator{}
	t.Logf("input %s", input)
	_, res := ki.attachContext([]byte(input), "321")

	dat := make(map[string]interface{})

	t.Logf("result %s, len: %d", string(res), len(res))

	if err := json.Unmarshal(res, &dat); err != nil {
		t.Errorf("on test error: %s", err)
	}

	if ctx := dat[CONTEXT]; ctx != nil {
		if cid := ctx.(map[string]interface{})[CORRELATION]; cid == nil {
			t.Errorf("missing correlation id")
		}
	} else {
		t.Errorf("missing context: %s", dat)
	}

}

func TestContextWithCorrelation(t *testing.T) {
	input := `
{
    "id": "A",
    "add": "B",
    "context": {
         "correlation": "dfb98431-c64b-0062-effa-9f16f52e5e31"
    }
}
`

	b := []byte(input)

	dat := make(map[string]interface{})

	if err := json.Unmarshal(b, &dat); err != nil {
		t.Errorf("on test error: %s", err)
	}

	if hasCid, cid := correlationInBody(b); !hasCid {
		t.Errorf("cid %s", cid)
	}

	ki := KafkaIntegrator{}
	_, res := ki.attachContext(b, "321")
	output := string(res)

	t.Logf("output %s, len: %d", output, len(output))

	if input != output {
		t.Errorf("not the same")
	}
}

func TestHasContextButWithoutCorrelation(t *testing.T) {
	input := `
{
    "id": "A",
    "add": "B",
    "context": {
         "some" :"someValue"
    }
}
`

	t.Logf("[before] output %s, len: %d", input, len(input))

	b := []byte(input)
	before := jsonAsMap(b, t)

	if !util.MatchNestedMapPath(before, []string{"context"}) {
		t.Errorf("[before][missing] root.context")
	}
	if !util.MatchNestedMapPath(before, []string{"context", "some"}) {
		t.Errorf("[before][missing] root.context.some")
	}
	if util.MatchNestedMapPath(before, []string{"context", "correlation"}) {
		t.Errorf("[before][has] root.context.correlation")
	}

	ki := KafkaIntegrator{}
	attached, res := ki.attachContext(b, "321")
	output := string(res)
	after := jsonAsMap(res, t)

	if !attached {
		t.Fail()
	}

	t.Logf("[after] output %s, len: %d", output, len(output))

	if !util.MatchNestedMapPath(after, []string{"context", "correlation"}) {
		t.Errorf("[after][missing] root.context.correlation")
	}
	if !util.MatchNestedMapPath(after, []string{"context", "some"}) {
		t.Errorf("[after][missing] root.context.some")
	}
}

func jsonAsMap(buff []byte, t *testing.T) map[string]interface{} {
	dat := make(map[string]interface{})

	if err := json.Unmarshal(buff, &dat); err != nil {
		t.Errorf("[json-parsing]: %s", err)
	}
	return dat
}
