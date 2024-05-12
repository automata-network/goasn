package types

import (
	"fmt"
	"net"

	"github.com/automata-network/goasn/embed"
	"github.com/automata-network/goasn/gen"
	"github.com/automata-network/goasn/utils"
)

type AsInfo struct {
	meta *gen.AsMeta
}

func NewAsInfo(meta *gen.AsMeta) AsInfo {
	return AsInfo{meta: meta}
}

func (a AsInfo) Name() string {
	return embed.ExtractString(gen.GEN_AS_NAME, a.meta.Name)
}

func (a AsInfo) Region() string {
	return embed.ExtractString(gen.GEN_REGION, a.meta.Region)
}

type CidrInfo struct {
	meta *gen.CidrMeta
}

func (c CidrInfo) GetArg(name string) interface{} {
	switch name {
	case "asn":
		return c.Asn()
	case "as.region":
		return c.AsRegion()
	case "as":
		return c.AsName()
	case "cidr":
		return c.Cidr()
	case "cidr.name":
		return c.Name()
	default:
		return nil
	}
}

type CidrCollection struct {
	Infos []CidrInfo
	Cnts  []*CidrCnt
}

type CidrCnt struct {
	Info CidrInfo
	Cnt  int
}

func (c *CidrCnt) GetArg(name string) interface{} {
	if name == "count" {
		return c.Cnt
	}
	return c.Info.GetArg(name)
}

func NewCidrInfo(meta *gen.CidrMeta) CidrInfo {
	return CidrInfo{meta}
}

func (c CidrInfo) Meta() *gen.CidrMeta {
	return c.meta
}

func (c CidrInfo) StartIP() net.IP {
	return utils.Int2IP(c.meta.Start)
}

func (c CidrInfo) EndIP() net.IP {
	return utils.Int2IP(c.meta.End)
}

func (c CidrInfo) Cidr() *net.IPNet {
	return utils.Range2CIDR(c.meta.Start, c.meta.End)
}

func (c CidrInfo) Name() string {
	return embed.ExtractString(gen.GEN_CIDR_NAME, c.meta.Name)
}

func (c CidrInfo) Asn() int {
	return int(c.meta.ASN)
}

func (c CidrInfo) AsName() string {
	ass := gen.AS_META_MAP[c.meta.ASN]
	info := AsInfo{ass[0]}
	return info.Name()
}

func (c CidrInfo) AsRegion() string {
	ass := gen.AS_META_MAP[c.meta.ASN]
	info := AsInfo{ass[0]}
	return info.Region()
}

func (c CidrInfo) String() string {
	return fmt.Sprintf("CidrInfo{Cidr: %v, ASN: %v, Name: %q}", c.Cidr(), c.Asn(), c.Name())
}
