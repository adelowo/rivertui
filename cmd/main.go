package main

import (
	"log"
	"os"

	"github.com/adelowo/rivertui/cmd/tui"
)

func main() {
	os.Setenv("TZ", "")

	if err := tui.Execute(); err != nil {
		log.Fatal(err)
	}
}
