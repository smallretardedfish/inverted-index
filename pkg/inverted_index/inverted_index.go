package inverted_index

import (
	"bufio"
	"fmt"
	set "github.com/smallretardedfish/inverted-index/pkg/hash_set"
	"golang.org/x/exp/maps"
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
	mu     sync.Mutex
	dict   map[string]set.HashSet[string]
}

type pair struct {
	word     string
	filename string
}

func (ii *MapInvertedIndex) Build(workersCount int, sources any) error {
	switch ii.source {
	case FileSourceType:
		filenames := sources.([]string)
		jobsCh := make(chan string, len(filenames))

		go func() { // feeling jobs channel with filenames
			for _, filename := range filenames {
				jobsCh <- filename
			}
			close(jobsCh)
		}()

		res := make(chan pair)
		go func() {
			wg := sync.WaitGroup{}
			for i := 0; i < workersCount; i++ { // worker pool
				wg.Add(1)
				go scanWords(&wg, jobsCh, res)
			}
			wg.Wait()
			close(res)
		}()

		ii.writeResult(res)

	case StringSourceType:
		stringSources := sources.([]StringSource)
		for _, src := range stringSources {
			words := strings.Split(src.Text, " ")
			for _, word := range words {
				if _, ok := ii.dict[word]; !ok {
					ii.dict[word] = make(set.HashSet[string])
				}
				ii.dict[word].Insert(src.Name)
			}
		}
	}

	workersStr := fmt.Sprintf("%d worker", workersCount)
	if workersCount > 1 {
		workersStr += "s"
	}
	log.Printf("inverted index was built by %s", workersStr)
	return nil
}

func (ii *MapInvertedIndex) writeResult(res <-chan pair) {
	for p := range res { // collecting the result to hashmap
		if _, ok := ii.dict[p.word]; !ok {
			ii.dict[p.word] = make(set.HashSet[string])
		}
		ii.dict[p.word].Insert(p.filename)
	}
}

func scanWords(wg *sync.WaitGroup, jobsCh <-chan string, ch chan<- pair) {
	defer wg.Done()
	for filename := range jobsCh { // job is  filename of file to process
		f, err := os.Open(filename)
		if err != nil {
			log.Println(err)
			return
		}

		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanWords)

		for scanner.Scan() {
			// TODO: check this
			word := strings.Trim(scanner.Text(), ".,/:';!@#$%&*()`~<>[]{}\n\r")
			ch <- pair{
				word:     word,
				filename: f.Name(),
			}
		}
		if err := f.Close(); err != nil {
			log.Println(err)
		}
	}
}

func (ii *MapInvertedIndex) Search(word string) []string {
	return maps.Keys(ii.dict[word])
}

func NewMapInvertedIndex(source IndexBuildSourceType) *MapInvertedIndex {
	return &MapInvertedIndex{
		source: source,
		dict:   make(map[string]set.HashSet[string]),
	}
}
