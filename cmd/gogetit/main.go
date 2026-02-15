package main

import (
	"os"

	"github.com/anxkhn/gogetit-workshop/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
