package util

import (
	"log"

	"github.com/oklog/ulid/v2"
	"s7i.io/kafka-gateway/internal/config"
)

const (
	KIND_UUID = iota + 10
	KIND_ULID
)

var kindOf = map[string]int{
	"uuid": KIND_UUID,
	"ulid": KIND_ULID,
}

type IdGen struct {
	kind int
}

func NewIdGen(app *config.CfgApp) *IdGen {
	k := kindOf[app.IdGenKind]
	if k == 0 {
		log.Fatalf("[Oops] Bad Id Gen Kind: %s, required: uuid or ulid\n", app.IdGenKind)
	}
	gen := IdGen{kind: k}
	return &gen
}

func (idGen *IdGen) MakeId() string {
	switch {
	case idGen.kind == KIND_UUID:
		return UUID()
	case idGen.kind == KIND_ULID:
		return ulid.Make().String()
	default:
		panic("oops")
	}
}
