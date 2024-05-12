package parser

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/automata-network/goasn/utils"
	"github.com/chzyer/logex"
)

type CidrInfo struct {
	Start uint32
	End   uint32
	ASN   int
	Name  string
}

func ParseCIDR(fp string) ([]*CidrInfo, error) {
	data, err := os.ReadFile(fp)
	if err != nil {
		return nil, logex.Trace(err)
	}
	gzipReader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, logex.Trace(err)
	}
	defer gzipReader.Close()

	r := bufio.NewReader(gzipReader)
	var cidrInfos []*CidrInfo
	for {
		lineBytes, err := r.ReadSlice('\n')
		if err != nil {
			break
		}
		lineSegs := bytes.SplitN(lineBytes, []byte("\t"), 2)
		ipAndCidr := bytes.SplitN(lineSegs[0], []byte("/"), 3)
		cidr := string(ipAndCidr[2])
		start, end, err := utils.CIDR2Range(cidr)
		if err != nil {
			return nil, logex.Trace(err)
		}
		asNo, err := strconv.Atoi(string(bytes.TrimPrefix(ipAndCidr[1], []byte("AS"))))
		if err != nil {
			return nil, logex.Trace(err)
		}
		cidrName := strings.ToLower(string(bytes.TrimSpace(lineSegs[1])))
		cidrInfos = append(cidrInfos, &CidrInfo{
			Start: start,
			End:   end,
			ASN:   asNo,
			Name:  cidrName,
		})
	}

	sort.Slice(cidrInfos, func(i, j int) bool {
		return cidrInfos[i].Start < cidrInfos[j].Start
	})

	return compactList(cidrInfos), nil
}

func compactList(list []*CidrInfo) []*CidrInfo {
	output := make([]*CidrInfo, 0, len(list))
	for i := range list {
		if i == 0 {
			output = append(output, list[0])
			continue
		}
		prev := output[len(output)-1]
		cur := list[i]
		if prev.End >= cur.Start {
			isLarger := prev.Start < cur.Start && prev.End > cur.End
			if isLarger {
				// merge
				continue
			}

			if prev.Start == cur.Start && prev.End == cur.End {
				// ignore
				continue
			}

			if prev.End == cur.End && prev.Start < cur.Start {
				// ignore
				continue

			}

			if prev.Start == cur.Start {
				if prev.End > cur.End {
					continue
				}
				if prev.End < cur.End {
					*prev = *cur
					continue
				}
			}

			fmt.Println("prev:", prev)
			fmt.Println("cur:", cur)
			panic("??")
		}
		output = append(output, cur)
	}
	return output
}
