package parser

import (
	"bytes"
	"encoding/binary"
	"os"

	"github.com/automata-network/goasn/utils"
	"github.com/chzyer/logex"
)

type StringInterning struct {
	lib   []string
	buf   []byte
	index map[string]int
}

func NewStringInterning(lib []string) *StringInterning {
	lib = utils.SortUniqueSlice(lib)

	buf := bytes.NewBuffer(nil)
	index := make(map[string]int, len(lib))
	for _, item := range lib {
		if len(item) > 65535 {
			panic(item)
		}
		index[item] = buf.Len()
		buf.Write(binary.BigEndian.AppendUint16(nil, (uint16(len(item)))))
		buf.WriteString(item)
	}
	return &StringInterning{
		lib:   lib,
		buf:   buf.Bytes(),
		index: index,
	}
}

func (si *StringInterning) GetIndex(s string) (int, bool) {
	idx, ok := si.index[s]
	return idx, ok
}

func (si *StringInterning) WriteFile(fp string) error {
	if err := os.WriteFile(fp, si.buf, 0755); err != nil {
		return logex.Trace(err)
	}
	return nil
}
