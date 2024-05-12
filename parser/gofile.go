package parser

import (
	"os"
	"os/exec"

	"github.com/chzyer/logex"
)

type GoFile struct {
}

func (g *GoFile) Write(fp string, data []byte) error {
	w, err := os.OpenFile(fp, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		return logex.Trace(err)
	}
	name := w.Name()
	defer w.Close()

	if _, err := w.Write(data); err != nil {
		return logex.Trace(err)
	}

	w.Close()

	cmd := exec.Command("go", "fmt", name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return logex.Trace(err)
	}

	return nil
}

func WriteGoFile(fp string, data []byte) error {
	return (&GoFile{}).Write(fp, data)
}
