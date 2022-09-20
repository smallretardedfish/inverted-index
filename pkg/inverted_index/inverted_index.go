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

type pair struct {
	word     string
	filename string
}

func (ii *MapInvertedIndex) Build(workersCount int, sources any) error {
	switch ii.source {
	case FileSourceType:
		filenames := sources.([]string)
		jobsCh := make(chan string, len(filenames))

		//chans := make([]chan pair, 0, workersCount)
		res := make(chan pair)
		wg := sync.WaitGroup{}
		wg.Add(workersCount)

		go func() { // feeling jobs channel with filenames
			for _, filename := range filenames {
				jobsCh <- filename
			}
			close(jobsCh)
		}()

		go func() {
			for p := range res { // collecting the result to hashmap
				if _, ok := ii.dict[p.word]; !ok {
					ii.dict[p.word] = make(hash_set.HashSet[string])
				}
				ii.dict[p.word].Insert(p.filename)
			}
		}()

		for i := 0; i < workersCount; i++ { // worker pool
			//pairCh := make(chan pair)     // every worker writes in this channel
			//chans = append(chans, pairCh) // TODO check alternatives of slice append here
			go func(id int) {
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
						word := scanner.Text()

						res <- pair{
							word:     word,
							filename: f.Name(),
						}
					}
				}
				//	log.Printf("worker %d done.\n", id+1)
			}(i)
		}
		wg.Wait()
		close(res)

		//res := utils.MergeChannels(chans...)

	case StringSourceType:
		stringSources := sources.([]StringSource)
		for _, src := range stringSources {
			words := strings.Split(src.Text, " ")
			for _, word := range words {
				if _, ok := ii.dict[word]; !ok {
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
