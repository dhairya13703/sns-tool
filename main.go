package main

import (
	"os"

	"github.com/dhairya13703/sns-tool/cmd/root"
)

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
