package main

import (
	"log"

	"github.com/RA341/bob/cli"
)

func main() {
	err := cli.NewApp().Run()
	if err != nil {
		log.Fatal(err)
	}
}
