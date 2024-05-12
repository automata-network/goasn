package gen

import (
	"github.com/automata-network/goasn/embed"
	"github.com/automata-network/goasn/gen/gen_embed"
)

var CIDR_META = embed.ExtractStructs[CidrMeta](gen_embed.GEN_CIDR_META)

var AS_META = embed.ExtractStructs[AsMeta](gen_embed.GEN_AS_META)

var AS_META_MAP = func() map[uint32][]*AsMeta {
	out := make(map[uint32][]*AsMeta, len(AS_META))
	for _, meta := range AS_META {
		out[meta.ASN] = append(out[meta.ASN], meta)
	}
	return out
}()

type AsMeta struct {
	Start  uint32
	End    uint32
	Region uint32
	Name   uint32
	ASN    uint32
}

type CidrMeta struct {
	Start uint32
	End   uint32
	Name  uint32
	ASN   uint32
}

func (a *AsMeta) GetStartIP() uint32 {
	return a.Start
}

func (a *AsMeta) GetEndIP() uint32 {
	return a.End
}

func (a *CidrMeta) GetStartIP() uint32 {
	return a.Start
}

func (a *CidrMeta) GetEndIP() uint32 {
	return a.End
}
