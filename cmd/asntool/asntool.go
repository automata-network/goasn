package main

import "github.com/chzyer/flagly"

type AsnTool struct {
	Generate *GenerateHandler `flagly:"handler"`
}

func main() {
	flagly.Run(AsnTool{})
}
