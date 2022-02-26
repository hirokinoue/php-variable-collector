package main

import (
	"io/ioutil"
	"path/filepath"
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

func filePaths(inDir, exclude string) ([]string, error) {
	files, err := ioutil.ReadDir(inDir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, f := range files {
		if f.Name() == exclude {
			continue
		}
		if f.IsDir() {
			tmpPaths, err := filePaths(filepath.Join(inDir, f.Name()), exclude)
			if err != nil {
				return nil, err
			}
			paths = append(paths, tmpPaths...)
			continue
		}
		paths = append(paths, filepath.Join(inDir, f.Name()))
	}
	return paths, nil
}
