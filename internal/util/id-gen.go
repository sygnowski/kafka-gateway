package util

import (
	"log"

	"github.com/oklog/ulid/v2"
	"s7i.io/kafka-gateway/internal/config"
)

type IdGen struct {
	kind string
}

func NewIdGen(app *config.CfgApp) *IdGen {
	gen := IdGen{kind: app.IdGenKind}
	return &gen
}

func (idGen *IdGen) MakeId() string {
	switch {
	case idGen.kind == "uuid":
		return UUID()
	case idGen.kind == "ulid":
		return ulid.Make().String()
	default:
		log.Fatalf("[Oops] bad Id Gen Kind: %s, required: uuid | ulid\n", idGen.kind)

		return ""
	}
}
