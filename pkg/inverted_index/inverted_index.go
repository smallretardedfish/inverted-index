package inverted_index

import (
	"bufio"
	"github.com/smallretardedfish/inverted-index/pkg/hash_set"
	"log"
	"os"
	"strings"
	"sync"
)

type InvertedIndex interface {
	Build(n int, sources any) error
	Search(word string) []string
}

type MapInvertedIndex struct {
	source IndexBuildSourceType
	mu     *sync.Mutex
	dict   map[string]hash_set.HashSet[string]
}

func (ii *MapInvertedIndex) Build(workersCount int, sources any) error {

	switch ii.source {
	case FileSourceType:
		filenames := sources.([]string)
		sourcesCh := make(chan int, len(filenames))

		go func() {
			for i := range filenames {
				sourcesCh <- i
			}
			close(sourcesCh)
		}()

		for i := 0; i < workersCount; i++ {
			go func() {
				for idx := range sourcesCh {
					f, err := os.Open(filenames[idx])
					if err != nil {
						log.Println(err)
						return
					}

					scanner := bufio.NewScanner(f)
					scanner.Split(bufio.ScanWords)

					for scanner.Scan() {
						word := scanner.Text()

						ii.mu.Lock()

						_, ok := ii.dict[word]
						if !ok {
							ii.dict[word] = make(hash_set.HashSet[string])
						}
						ii.dict[word].Insert(f.Name())

						ii.mu.Unlock()
					}
				}
			}()
		}
	//
	case StringSourceType:
		stringSources := sources.([]StringSource)
		for _, src := range stringSources {
			words := strings.Split(src.Text, " ")

			for _, word := range words {
				_, ok := ii.dict[word]
				if !ok {
					ii.dict[word] = make(hash_set.HashSet[string])
				}

				ii.dict[word].Insert(src.Name)
			}
		}
	}

	return nil
}

func (ii *MapInvertedIndex) Search(word string) []string {
	set := ii.dict[word]
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
