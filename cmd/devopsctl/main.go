package main

import (
	"devopsctl/internal/devopsctlcmd"
	"os"
)

func main() {
	code := devopsctlcmd.Main()
	os.Exit(int(code))
}
