package main

import (
	"strings"
)

func isPhpFile(s string) bool {
	return strings.Contains(s, ".php")
}
