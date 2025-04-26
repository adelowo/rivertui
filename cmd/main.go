package main

import (
	"log"
	"os"

	"github.com/adelowo/rivertui/cmd/cli"
)

func main() {
	os.Setenv("TZ", "")

	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
