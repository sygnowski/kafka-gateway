package kafint

import (
	"encoding/json"
	"testing"
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

	ki := KafkaIntegrator{}
	_, res := ki.attachContext(b, "321")
	output := string(res)

	t.Logf("output %s, len: %d", output, len(output))

	if input != output {
		t.Errorf("not the same")
	}
}
