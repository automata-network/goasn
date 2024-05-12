package parser

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/chzyer/logex"
)

type IpCnt struct {
	IP  net.IP
	Cnt int
}

func ParseIpCnt(r io.Reader) ([]*IpCnt, error) {
	buf := bufio.NewReader(r)
	idx := 0
	var out []*IpCnt
	for {
		lineBytes, err := buf.ReadSlice('\n')
		if err != nil {
			break
		}
		line := strings.Split(string(bytes.TrimSpace(lineBytes)), " ")

		ip := net.ParseIP(line[len(line)-1])
		if ip == nil {
			return nil, logex.Trace(err, string(lineBytes))
		}
		count := 1
		if len(line) > 1 {
			count, _ = strconv.Atoi(line[0])
			if count == 0 {
				count = 1
			}
		}
		idx++
		out = append(out, &IpCnt{
			IP:  ip.To4(),
			Cnt: count,
		})
	}
	return out, nil
}
