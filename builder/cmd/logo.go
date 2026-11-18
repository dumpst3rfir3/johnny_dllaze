package main

import (
	_ "embed"
	"fmt"
)

//go:embed logo.ansi
var logoData string

func printLogo() {
	fmt.Print(logoData)
}
