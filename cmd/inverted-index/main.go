package main

import (
	"bufio"
	"fmt"
	"github.com/smallretardedfish/inverted-index/pkg/inverted_index"
	"log"
	"os"
	"strings"
)

const (
	FirstText = `In computer science, an inverted index (also referred to as a postings list, postings file,
or inverted file) is a database index storing a mapping from content, such as words or numbers,
to its locations in a table, or in a document or a set of documents
(named in contrast to a forward index, which maps from documents to content).
The purpose of an inverted index is to allow fast full-text searches,
at a cost of increased processing when a document is added to the database.
The inverted file may be the database file itself, rather than its index. 
It is the most popular data structure used in document retrieval systems`
	SecondText = "THIS IS computer science inverted index"
)

func run() error {
	sources := []inverted_index.StringSource{{
		Name: "First",
		Text: FirstText,
	}, {
		Name: "Second",
		Text: SecondText,
	}}

	invInvex := inverted_index.NewMapInvertedIndex(inverted_index.StringSourceType)

	if err := invInvex.Build(sources); err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		word := scanner.Text()
		res := invInvex.Search(word)
		fmt.Printf("word: %s, sources: %s\n", word, strings.Join(res, ","))
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}
