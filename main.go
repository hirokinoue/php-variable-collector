package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
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

func phpFilePaths(inDir, exclude string) ([]string, error) {
	paths, err := filePaths(inDir, exclude)
	if err != nil {
		return nil, err
	}
	var phps []string
	for _, p := range paths {
		if !isPhpFile(p) {
			continue
		}
		phps = append(phps, p)
	}
	return phps, nil
}

type dict struct {
	value map[string]bool
	mux   sync.Mutex
}

func newDict() *dict {
	return &dict{
		value: make(map[string]bool),
	}
}

func (d *dict) add(variable string) {
	d.mux.Lock()
	if _, ok := d.value[variable]; !ok {
		d.value[variable] = true
	}
	d.mux.Unlock()
}

func (d *dict) sortValue() []string {
	d.mux.Lock()
	defer d.mux.Unlock()
	keys := make([]string, 0, len(d.value))
	for k := range d.value {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func collectPhpVariable(filePath string, ch chan<- []string, e chan<- error, semaphore chan struct{}) {
	semaphore <- struct{}{}
	defer func() {
		<-semaphore
	}()
	f, err := os.Open(filePath)
	if err != nil {
		e <- err
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	var strs []string
	for scanner.Scan() {
		words := strings.Split(scanner.Text(), " ")
		for _, w := range words {
			if isPhpVariable(w) {
				strs = append(strs, removeSymbolFromVariable(w))
			}
		}
	}
	ch <- strs
}

func writeFile(outFilePath string, line string) error {
	file, err1 := os.OpenFile(outFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err1 != nil {
		return err1
	}
	defer file.Close()

	_, err2 := file.WriteString(fmt.Sprintf("%s%s", line, "\n"))
	if err2 != nil {
		return err2
	}
	return nil
}
