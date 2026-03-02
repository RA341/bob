package cli

import (
	"fmt"
	"os"
	"path/filepath"
)

type CmdInit struct {
	flags []*Flag
}

func (c *CmdInit) Name() string {
	return "init"
}

func (c *CmdInit) Help() string {
	return "Creates a sample Bobfile in current working directory, pass optional --dir <path> for create in a specific path"
}

func (c *CmdInit) Run(args ...string) error {
	dirFlag := StrFlag("dir", "Set custom dir for flag")

	_, err := ParseFlags([]*Flag{
		dirFlag,
	}, args)
	if err != nil {
		return err
	}

	var path string
	if dirFlag.isSet {
		var ok bool
		path, ok = dirFlag.val.(string)
		if !ok {
			return fmt.Errorf("flag `dir` must be a string")
		}
	} else {
		path, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	fPath := filepath.Join(path, BobFileName)
	fmt.Println("Creating file at", fPath)

	file, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	_, err = file.WriteString(`hello(){
	print hello world
}`)

	return err
}
