package main

import (
	"fmt"
	"sort"
	"strings"
)

// ====== Document ======
type Document struct {
	ID   int
	Text string
}

// ====== Indexer ======
type Indexer interface {
	Index(docs []Document)
	Search(keyword string) []int
}

type InvertedIndexer struct {
	index map[string][]int
}

func NewInvertedIndexer() *InvertedIndexer {
	return &InvertedIndexer{index: make(map[string][]int)}
}

func (i *InvertedIndexer) Index(docs []Document) {
	for _, doc := range docs {
		words := strings.Fields(strings.ToLower(doc.Text))
		seen := make(map[string]bool)
		for _, word := range words {
			if !seen[word] {
				i.index[word] = append(i.index[word], doc.ID)
				seen[word] = true
			}
		}
	}
}

func (i *InvertedIndexer) Search(keyword string) []int {
	return i.index[strings.ToLower(keyword)]
}

// ====== Ranking Strategy Pattern ======
type RankingStrategy interface {
	Rank(results []int, docs map[int]Document, keyword string) []int
}

type ByDocSize struct{}

func (r *ByDocSize) Rank(results []int, docs map[int]Document, keyword string) []int {
	sort.Slice(results, func(i, j int) bool {
		return len(docs[results[i]].Text) < len(docs[results[j]].Text)
	})
	return results
}

type ByKeywordFrequency struct{}

func (r *ByKeywordFrequency) Rank(results []int, docs map[int]Document, keyword string) []int {
	keyword = strings.ToLower(keyword)
	sort.Slice(results, func(i, j int) bool {
		return strings.Count(strings.ToLower(docs[results[i]].Text), keyword) >
			strings.Count(strings.ToLower(docs[results[j]].Text), keyword)
	})
	return results
}

// ====== Ranking Factory ======
func GetRankingStrategy(method string) RankingStrategy {
	switch method {
	case "size":
		return &ByDocSize{}
	case "frequency":
		return &ByKeywordFrequency{}
	default:
		return &ByDocSize{} // Default strategy
	}
}

// ====== Search Engine ======
type SearchEngine struct {
	documents map[int]Document
	indexer   Indexer
}

func NewSearchEngine(indexer Indexer) *SearchEngine {
	return &SearchEngine{
		documents: make(map[int]Document),
		indexer:   indexer,
	}
}

func (s *SearchEngine) AddDocuments(docs []Document) {
	for _, doc := range docs {
		s.documents[doc.ID] = doc
	}
	s.indexer.Index(docs)
}

func (s *SearchEngine) Search(keyword, rankingMethod string) []Document {
	ids := s.indexer.Search(keyword)
	ranker := GetRankingStrategy(rankingMethod)
	sortedIDs := ranker.Rank(ids, s.documents, keyword)

	results := make([]Document, 0, len(sortedIDs))
	for _, id := range sortedIDs {
		results = append(results, s.documents[id])
	}
	return results
}

// ====== Main ======
func main() {
	docs := []Document{
		{ID: 1, Text: "Go is expressive, concise, clean, and efficient."},
		{ID: 2, Text: "Concurrency is not parallelism."},
		{ID: 3, Text: "Go makes it easy to build simple, reliable, and efficient software."},
	}

	searchEngine := NewSearchEngine(NewInvertedIndexer())
	searchEngine.AddDocuments(docs)

	fmt.Println("Search by frequency:")
	results := searchEngine.Search("efficient", "frequency")
	for _, doc := range results {
		fmt.Printf("Doc %d: %s\n", doc.ID, doc.Text)
	}

	fmt.Println("\nSearch by document size:")
	results = searchEngine.Search("efficient", "size")
	for _, doc := range results {
		fmt.Printf("Doc %d: %s\n", doc.ID, doc.Text)
	}
}
