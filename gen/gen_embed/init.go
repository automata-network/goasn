package gen_embed

import (
	_ "embed"
)

//go:embed gen_as_meta.embed
var GEN_AS_META []byte

//go:embed gen_cidr_name.embed
var GEN_CIDR_NAME []byte

//go:embed gen_cidr_meta.embed
var GEN_CIDR_META []byte

//go:embed gen_region.embed
var GEN_REGION []byte

//go:embed gen_as_name.embed
var GEN_AS_NAME []byte
