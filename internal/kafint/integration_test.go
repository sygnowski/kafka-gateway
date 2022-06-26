package kafint

import (
	"encoding/json"
	"testing"
)

func Test_attachContext(t *testing.T) {
	input := `{"id": "123"}`

	ki := KafkaIntegrator{}
	t.Logf("input %s", input)
	res := ki.attachContext([]byte(input), "321")

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
