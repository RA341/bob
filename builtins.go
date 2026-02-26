package main

import (
	"fmt"
	"os"
)

type BuiltInFunc func(vals ...string) error

var builtInMaps = map[string]BuiltInFunc{
	"print": func(vals ...string) error {
		var interfaceArgs []interface{}
		for _, val := range vals {
			interfaceArgs = append(interfaceArgs, val)
		}

		fmt.Println(interfaceArgs...)
		return nil
	},
	"workdir": func(vals ...string) error {
		return os.Chdir(vals[0])
	},
}
