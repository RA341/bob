package cli

import (
	"os"
	"path/filepath"
)

type CmdInit struct{}

func (c *CmdInit) Name() string {
	return "init"
}

func (c *CmdInit) Help() string {
	return "Creates a sample Bobfile in current working directory"
}

func (c *CmdInit) Run() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	fPath := filepath.Join(wd, BobFileName)
	file, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(`hello(){
	print hello world
}`)

	return err
}
