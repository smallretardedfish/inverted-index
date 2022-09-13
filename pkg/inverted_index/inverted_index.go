package inverted_index

import (
	"bufio"
	"github.com/smallretardedfish/inverted-index/pkg/hash_set"
	"os"
	"strings"
	"sync"
)

type InvertedIndex interface {
	Build(sources any) error
	Search(word string) []string
}

type MapInvertedIndex struct {
	source IndexBuildSourceType
	mu     *sync.Mutex
	dict   map[string]hash_set.HashSet[string]
}

func (i *MapInvertedIndex) Build(sources any) error {
	switch i.source {
	case FileSourceType:
		filenames := sources.([]string)
		for _, filename := range filenames {
			f, err := os.Open(filename)
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(f)
			scanner.Split(bufio.ScanWords)

			for scanner.Scan() {
				word := scanner.Text()
				_, ok := i.dict[word]
				if !ok {
					i.dict[word] = make(hash_set.HashSet[string])
				}
				i.dict[word].Insert(f.Name())
			}
		}

	case StringSourceType:
		stringSources := sources.([]StringSource)
		for _, src := range stringSources {
			words := strings.Split(src.Text, " ")

			for _, word := range words {
				_, ok := i.dict[word]
				if !ok {
					i.dict[word] = make(hash_set.HashSet[string])
				}

				i.dict[word].Insert(src.Name)
			}
		}
	}

	return nil
}

func (i *MapInvertedIndex) Search(word string) []string {
	set := i.dict[word]
	res := make([]string, 0, set.Size())

	for source := range set {
		res = append(res, source)
	}

	return res
}

func NewMapInvertedIndex(source IndexBuildSourceType) *MapInvertedIndex {
	return &MapInvertedIndex{
		source: source,
		mu:     &sync.Mutex{},
		dict:   make(map[string]hash_set.HashSet[string]),
	}
}
