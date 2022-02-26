package main

import (
	"strings"
)

func isPhpFile(s string) bool {
	return strings.Contains(s, ".php")
}

func isPhpVariable(s string) bool {
	if strings.Index(s, "$") != 0 {
		return false
	}
	if strings.Contains(s, "->") {
		return false
	}
	return true
}
