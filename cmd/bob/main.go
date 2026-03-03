package main

import (
	"log"
	"os"

	"github.com/RA341/bob/cli"
)

func main() {
	err := cli.NewApp().Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
