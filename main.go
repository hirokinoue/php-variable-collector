package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
)

const outputFileName = "variables.txt"

func main() {
	log.SetFlags(log.Llongfile)
	inDir, outDir, exclude := parseFlags(os.Args[1:])
	err := clearOutDir(outDir)
	if err != nil {
		log.Println(err)
	}

	d := createVariableDictionary(inDir, exclude)

	for _, v := range d.sortValue() {
		err := writeFile(filepath.Join(outDir, outputFileName), v)
		if err != nil {
			log.Println(err)
		}
	}
	fmt.Println("Done")
}

func createVariableDictionary(inDir, exclude string) *dict {
	d := newDict()
	paths, err := phpFilePaths(inDir, exclude)
	if err != nil {
		log.Println(err)
	}
	ch := make(chan []string, len(paths))
	e := make(chan error)
	semaphore := make(chan struct{}, runtime.NumCPU())
	for _, p := range paths {
		go collectPhpVariable(p, ch, e, semaphore)
	}
	for i := 0; i < len(paths); i++ {
		select {
		case strs := <-ch:
			for _, s := range strs {
				d.add(s)
			}
		case err := <-e:
			log.Println(err)
		}
	}
	close(ch)
	close(e)
	close(semaphore)
	return d
}

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
	sub := regexp.MustCompile(`[\[\].,;!"')(:%+-]`).Split(s, -1)
	return sub[0]
}

func filePaths(inDir, exclude string) ([]string, error) {
	files, err := os.ReadDir(inDir)
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

func clearOutDir(outDir string) error {
	paths, err := filePaths(outDir, "")
	if err != nil {
		return err
	}
	for _, p := range paths {
		err := os.Remove(p)
		if err != nil {
			return err
		}
	}
	return nil
}

func parseFlags(args []string) (string, string, string) {
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	inDir := f.String("in", "in", "入力ファイルを格納するディレクトリ")
	outDir := f.String("out", "out", "処理結果が出力されるディレクトリ")
	exclude := f.String("exclude", "", "処理しないファイルやディレクトリ")
	// ExitOnErrorモードなので戻り値がerrorになることはない
	_ = f.Parse(args)
	return *inDir, *outDir, *exclude
}
