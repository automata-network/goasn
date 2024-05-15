package goasn

import (
	"net"
	"sort"

	"github.com/automata-network/goasn/gen"
	"github.com/automata-network/goasn/parser"
	"github.com/automata-network/goasn/types"
	"github.com/automata-network/goasn/utils"
)

func GenerateAllCidrInfo() []types.CidrInfo {
	out := make([]types.CidrInfo, len(gen.CIDR_META))
	for idx, meta := range gen.CIDR_META {
		out[idx] = types.NewCidrInfo(meta)
	}
	return out
}

func SearchCidr(ip net.IP) (types.CidrInfo, bool) {
	meta, ok := utils.FindIpInfo(gen.CIDR_META, ip.To4())
	if !ok {
		return types.NewCidrInfo(nil), false
	}
	return types.NewCidrInfo(meta), true
}

func CollectCidrsByIpCnt(cnts []*parser.IpCnt) types.CidrCollection {
	totalCount := make(map[*gen.CidrMeta]int)
	for _, ipcnt := range cnts {
		cidr, ok := SearchCidr(ipcnt.IP)
		if !ok {
			continue
		}

		totalCount[cidr.Meta()] = totalCount[cidr.Meta()] + ipcnt.Cnt
	}

	counts := make([]*types.CidrCnt, 0, len(totalCount))
	infos := make([]types.CidrInfo, 0, len(totalCount))

	for meta, cnt := range totalCount {
		counts = append(counts, &types.CidrCnt{
			Info: types.NewCidrInfo(meta),
			Cnt:  cnt,
		})
	}
	sort.Slice(counts, func(i, j int) bool {
		return counts[i].Cnt < counts[j].Cnt
	})
	for _, item := range counts {
		infos = append(infos, item.Info)
	}
	return types.CidrCollection{
		Infos: infos,
		Cnts:  counts,
	}
}
