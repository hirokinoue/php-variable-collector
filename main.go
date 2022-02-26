package main

import (
	"regexp"
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

func removeSymbolFromVariable(s string) string {
	sub := regexp.MustCompile(`[\[\]\.,;!"')(:%+-]`).Split(s, -1)
	return sub[0]
}
