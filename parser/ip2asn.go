package parser

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"math"
	"net"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/automata-network/goasn/utils"
	"github.com/chzyer/logex"
)

type AsInfo struct {
	Start  net.IP
	End    net.IP
	ASN    int
	Region string
	Name   string
}

func (i *AsInfo) TryMerge(other *AsInfo) error {
	if i.ASN != other.ASN {
		return logex.NewError("asn not match")
	}
	if i.Region != other.Region {
		return logex.NewError("region not match")
	}
	if i.Name != other.Name {
		return logex.NewError("name not match")
	}
	if utils.IP2Int(i.End)+1 == utils.IP2Int(other.Start) || utils.IP2Int(other.End)+1 == utils.IP2Int(other.Start) {
		start := min(utils.IP2Int(other.Start), utils.IP2Int(i.Start))
		end := max(utils.IP2Int(other.End), utils.IP2Int(i.End))
		i.Start = utils.Int2IP(start)
		i.End = utils.Int2IP(end)
		return nil
	}
	return logex.NewErrorf("can't merge: origin: %v, target: %v", i, other)
}

func ParseIP2asn(fp string) ([]*AsInfo, error) {
	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, logex.Trace(err)
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, logex.Trace(err)
	}
	r := bufio.NewReader(gzipReader)
	total := 0
	infos := make([]*AsInfo, 0, 1024000)
	for {
		total += 1
		line, err := r.ReadSlice('\n')
		if err != nil {
			break
		}
		sp := bytes.SplitN(line, []byte("\t"), 5)
		start := net.ParseIP(string(sp[0])).To4()
		end := net.ParseIP(string(sp[1])).To4()
		asNo, err := strconv.Atoi(string(sp[2]))
		if err != nil {
			return nil, logex.Trace(err)
		}
		if asNo == 0 {
			continue
		}
		if asNo > int(math.MaxUint32) {
			panic("invalid asno")
		}
		asRegion := strings.TrimSpace(string(sp[3]))
		asName := strings.ToLower(strings.TrimSpace(string(sp[4])))

		infos = append(infos, &AsInfo{
			Start: start, End: end, ASN: asNo, Region: asRegion, Name: asName,
		})
	}

	// asnIndex := make(map[int]int, len(infos))
	// for i := 0; i < len(infos); i++ {
	// 	if _, ok := asnIndex[infos[i].ASN]; ok {
	// 		current := infos[i]
	// 		old := infos[asnIndex[infos[i].ASN]]
	// 		if err := old.TryMerge(current); err != nil {
	// 			return nil, logex.Trace(err)
	// 		}
	// 	} else {
	// 		asnIndex[infos[i].ASN] = i
	// 	}
	// }

	sort.Slice(infos, func(i, j int) bool {
		return utils.IP2Int(infos[i].Start) < utils.IP2Int(infos[j].Start)
	})
	return infos, nil
}
