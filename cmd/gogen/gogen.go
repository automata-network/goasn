package main

import (
	"os"

	"github.com/automata-network/goasn/embed"
	"github.com/automata-network/goasn/gen"
	"github.com/automata-network/goasn/parser"
	"github.com/automata-network/goasn/utils"
	"github.com/chzyer/logex"
)

func main() {
	if err := generateIP2asn(); err != nil {
		logex.Fatal(err)
	}
}

func generateIP2asn() error {
	asInfos, err := parser.ParseIP2asn("assets/ip2asn-v4.tsv.gz")
	if err != nil {
		return logex.Trace(err)
	}

	// REGIONS
	regions := parser.NewStringInterning(utils.CollectSliceValues(asInfos, func(n *parser.AsInfo) string {
		return n.Region
	}))
	if err := regions.WriteFile("gen/gen_embed/gen_region.embed"); err != nil {
		return logex.Trace(err)
	}

	// ASNAMES
	asNames := parser.NewStringInterning(utils.CollectSliceValues(asInfos, func(n *parser.AsInfo) string {
		return n.Name
	}))
	if err := asNames.WriteFile("gen/gen_embed/gen_as_name.embed"); err != nil {
		return logex.Trace(err)
	}

	if err := generateAsMeta(asInfos, regions, asNames); err != nil {
		return logex.Trace(err)
	}

	cidrs, err := parser.ParseCIDR("assets/as_cidrs.txt.gz")
	if err != nil {
		return logex.Trace(err)
	}
	cidrNames := parser.NewStringInterning(utils.CollectSliceValues(cidrs, func(n *parser.CidrInfo) string {
		return n.Name
	}))
	if err := cidrNames.WriteFile("gen/gen_embed/gen_cidr_name.embed"); err != nil {
		return logex.Trace(err)
	}
	if err := generateCidrMeta(cidrs, cidrNames); err != nil {
		return logex.Trace(err)
	}
	return nil
}

func generateCidrMeta(cidrs []*parser.CidrInfo, cidrNames *parser.StringInterning) error {
	metas := make([]*gen.CidrMeta, len(cidrs))
	for idx, cidr := range cidrs {
		name, ok := cidrNames.GetIndex(cidr.Name)
		if !ok {
			panic(cidr.Name + "not found")
		}
		metas[idx] = &gen.CidrMeta{
			Start: cidr.Start,
			End:   cidr.End,
			Name:  uint32(name),
			ASN:   uint32(cidr.ASN),
		}
	}
	if err := os.WriteFile("gen/gen_embed/gen_cidr_meta.embed", embed.EmbedStructs(metas), 0755); err != nil {
		return logex.Trace(err)
	}
	return nil
}

func generateAsMeta(asInfos []*parser.AsInfo, regions, asNames *parser.StringInterning) error {
	out := make([]*gen.AsMeta, len(asInfos))
	for idx, as := range asInfos {
		region, ok := regions.GetIndex(as.Region)
		if !ok {
			logex.Fatal(as.Region, "not found")
		}
		name, ok := asNames.GetIndex(as.Name)
		if !ok {
			panic(as.Name + "not found")
		}
		out[idx] = &gen.AsMeta{
			Start:  utils.IP2Int(as.Start),
			End:    utils.IP2Int(as.End),
			Region: uint32(region),
			Name:   uint32(name),
			ASN:    uint32(as.ASN),
		}
	}

	if err := os.WriteFile("gen/gen_embed/gen_as_meta.embed", embed.EmbedStructs(out), 0755); err != nil {
		return logex.Trace(err)
	}
	return nil
}
