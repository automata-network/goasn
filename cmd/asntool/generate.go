package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/automata-network/goasn"
	"github.com/automata-network/goasn/parser"
	"github.com/chzyer/logex"
)

type GenerateHandler struct {
	Format     string
	Filter     string
	FilterFile string
	Invert     bool `name:"v"`
}

func (g *GenerateHandler) FlaglyHandle() error {
	var filters []string
	if len(g.Filter) > 0 {
		for _, item := range strings.Split(g.Filter, ";") {
			filters = append(filters, strings.TrimSpace(item))
		}
	}
	if g.FilterFile != "" {
		data, err := os.ReadFile(g.FilterFile)
		if err != nil {
			return logex.Trace(err)
		}
		filters = strings.Split(string(data), "\n")
	}

	rules := goasn.ParseRules(filters)

	ipcnts, err := parser.ParseIpCnt(os.Stdin)
	if err != nil {
		return logex.Trace(err)
	}

	cidrCol := goasn.CollectCidrsByIpCnt(ipcnts)

	ruleCidrs := rules.FilterCidrs(cidrCol.Infos)
	fi, err := goasn.NewFormatInfo(g.Format)
	if err != nil {
		return logex.Trace(err)
	}
	total := 0
	for _, info := range cidrCol.Cnts {
		if len(ruleCidrs) > 0 {
			isExist := false
			for _, rc := range ruleCidrs {
				if info.Info == rc.Cidr {
					isExist = true
					break
				}
			}
			if g.Invert {
				isExist = !isExist
			}
			if !isExist {
				continue
			}
		}
		total += info.Cnt
		fmt.Println(fi.Format(info))
	}
	fmt.Println("total:", total)

	return nil
}
