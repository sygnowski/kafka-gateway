package util

import (
	"testing"

	"s7i.io/kafka-gateway/internal/config"
)

func Test_IdGen_UUID(t *testing.T) {
	gen := NewIdGen(&config.CfgApp{IdGenKind: "uuid"})
	id := gen.MakeId()

	t.Logf("UUID: %s", id)

}

func Test_IdGen_ULID(t *testing.T) {
	gen := NewIdGen(&config.CfgApp{IdGenKind: "ulid"})
	id := gen.MakeId()

	t.Logf("LUID: %s", id)

}
