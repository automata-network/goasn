package goasn

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/chzyer/logex"
)

type FormatInfo struct {
	layout    string
	variables []string
	isNeedAsn bool
}

func NewFormatInfo(format string) (*FormatInfo, error) {
	newFormat, err := strconv.Unquote(`"` + format + `"`)
	if err != nil {
		return nil, logex.Trace(err, format)
	}
	format = newFormat
	varsRaw := strings.Split(format, "@")
	var vars []string
	var newLayout []string

	isNeedAsn := false
	for idx, item := range varsRaw {
		nameEndIdx := strings.IndexAny(item, " \t/=|")
		name := item
		if nameEndIdx >= 0 {
			name = item[:nameEndIdx]
		}
		if idx > 0 {
			if nameEndIdx >= 0 {
				item = "%v" + item[nameEndIdx:]
			} else {
				item = "%v"
			}
			vars = append(vars, name)
		}
		newLayout = append(newLayout, item)
		if name == "as" || name == "asn.region" {
			isNeedAsn = true
		}
	}
	layout := strings.Join(newLayout, "")
	return &FormatInfo{
		variables: vars,
		layout:    layout,
		isNeedAsn: isNeedAsn,
	}, nil
}

type FormatArgs interface {
	GetArg(name string) interface{}
}

func (f *FormatInfo) Format(namedArgs FormatArgs) string {
	args := make([]interface{}, len(f.variables))
	for idx, vName := range f.variables {
		args[idx] = namedArgs.GetArg(vName)
	}
	return fmt.Sprintf(f.layout, args...)
}
