package goasn

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/automata-network/goasn/types"
	"github.com/chzyer/logex"
)

type Rules []*Rule

type RuleCidrPair struct {
	Rule *Rule
	Cidr types.CidrInfo
}

func (rules Rules) FilterCidrs(cidrList []types.CidrInfo) []*RuleCidrPair {
	var list []*RuleCidrPair
nextLoop:
	for _, c := range cidrList {
		for _, rule := range rules {
			if rule.MatchCIDR(c) {
				list = append(list, &RuleCidrPair{rule, c})
				continue nextLoop
			}
		}
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].Cidr.Meta().Start < list[j].Cidr.Meta().Start
	})
	return list

}

type Rule struct {
	Name    string
	ASN     int
	AS      string
	Region  string
	CIDR    string
	Comment string
}

func (rule *Rule) MatchCIDR(info types.CidrInfo) bool {
	if rule.Name != "" {
		if strings.Contains(info.Name(), rule.Name) {
			return true
		}
	}
	if rule.ASN > 0 {
		if info.Asn() == rule.ASN {
			return true
		}
	}
	if rule.AS != "" {
		if strings.Contains(info.AsName(), rule.AS) {
			return true
		}
	}
	if rule.Region != "" {
		if info.AsRegion() == rule.Region {
			return true
		}
	}
	if rule.CIDR != "" {
		if info.Cidr().String() == rule.CIDR {
			return true
		}
	}
	return false
}

func (r *Rule) String() string {
	output := ""
	if r.Name != "" {
		output = fmt.Sprintf("NAME:%v", r.Name)
	} else if r.AS != "" {
		output = fmt.Sprintf("AS:%v", r.AS)
	} else if r.ASN > 0 {
		output = fmt.Sprintf("ASN:%v", r.ASN)
	} else if r.Region != "" {
		output = fmt.Sprintf("REGION:%v", r.Region)
	} else if r.CIDR != "" {
		output = fmt.Sprintf("CIDR:%v", r.CIDR)
	} else {
		output = "{Rule}"
	}
	// if r.Comment != "" {
	// output += " #" + r.Comment
	// }
	return output
}

func ParseRules(lines []string) Rules {
	rules := make([]*Rule, 0, len(lines))
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		sp := strings.SplitN(line, ":", 2)
		ruleName := sp[0]
		ruleValues := strings.SplitN(sp[1], "#", 2)
		ruleValue := strings.TrimSpace(ruleValues[0])
		rule := &Rule{}
		if len(ruleValues) > 1 {
			rule.Comment = ruleValues[1]
		}

		switch ruleName {
		case "NAME":
			rule.Name = strings.ToLower(ruleValue)
		case "AS":
			rule.AS = strings.ToLower(ruleValue)
		case "ASN":
			asn, err := strconv.Atoi(ruleValue)
			if err != nil {
				logex.Fatal(err)
			}
			rule.ASN = asn
		case "CIDR":
			rule.CIDR = strings.ToLower(ruleValue)
		case "REGION":
			rule.Region = ruleValue
		default:
			continue
		}
		rules = append(rules, rule)
	}
	return rules
}
