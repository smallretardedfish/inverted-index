package main

import (
	"bufio"
	"fmt"
	utils "github.com/smallretardedfish/inverted-index/pkg"
	"github.com/smallretardedfish/inverted-index/pkg/inverted_index"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func run() error {
	args := make([]string, 2)
	copy(args, os.Args[1:])

	if len(args) == 0 {
		log.Println("no args given")
		args = append(args, "1", ".")
	} else if len(args) == 1 {
		args = append(args, ".")
	}

	num := args[0]
	n, err := strconv.Atoi(num)
	if err != nil {
		return fmt.Errorf("number of workers must be integer: %w", err)
	}

	dir := args[1]
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("can't read from given dir: %s, error: %w", dir, err)
	}

	var fileSources []string
	for _, entry := range dirEntries {
		fileSources = append(fileSources, filepath.Join(dir, entry.Name()))
	}

	fmt.Printf("number of workers:%d\ndirectory with files to be processed:%s\n", n, dir)
	invIndex := inverted_index.NewMapInvertedIndex(inverted_index.FileSourceType)

	var (
		e error
	)

	t := utils.EstimateExecutionTime(func() {
		if err := invIndex.Build(n, fileSources); err != nil {
			e = err
		}
	})
	if e != nil {
		return e
	}
	log.Println(t)

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf(">> ")
	for scanner.Scan() {
		word := scanner.Text()
		res := invIndex.Search(word)
		fmt.Printf("word: '%s' \n sources: %s\n>> ", word, strings.Join(res, ","))
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
