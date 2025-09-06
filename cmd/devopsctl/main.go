package main

import (
	"devopsctl/internal/cmd"
	"os"
)

func main() {
	code := cmd.Main()
	os.Exit(int(code))
}
