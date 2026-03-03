package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RA341/bob/cli/bob"
	"github.com/RA341/bob/cli/flag"
)

type Init struct {
	flags []*flag.Flag
}

func (c *Init) Name() string {
	return "init"
}

func (c *Init) Help() string {
	helpText := "Creates a sample Bobfile in current working directory"
	return PrintHelpSubCmd(c.Name(), c.flags, helpText)
}

func (c *Init) LoadFlags(args []string) error {
	c.flags = []*flag.Flag{
		flag.StrFlag("dir", "Set custom path for Bobfile (Optional)"),
		flag.HelpFlag(""),
	}
	_, err := flag.ParseFlags(c.flags, args)
	return err
}

func (c *Init) Run() (err error) {
	if c.flags[1].IsSet {
		fmt.Print(c.Help())
		return nil
	}

	var path string
	dirFlag := c.flags[0]
	if dirFlag.IsSet {
		var ok bool
		path, ok = dirFlag.Val.(string)
		if !ok {
			return fmt.Errorf("flag `dir` must be a string")
		}
	} else {
		path, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	fPath := filepath.Join(path, bob.FileName)
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

	_, err = file.WriteString("hello(){\nprint hello world\n}")
	return err
}
