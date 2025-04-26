package main

import (
	"fmt"
	"sort"
	"strings"
)

// ====== Document ======
type Document struct {
	ID       int
	Text     string
	Category string
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

// ====== Category Indexer (Keyword-style) ======
type CategoryIndexer struct {
	categoryIndex map[string]map[int]struct{}
}

func NewCategoryIndexer() *CategoryIndexer {
	return &CategoryIndexer{
		categoryIndex: make(map[string]map[int]struct{}),
	}
}

func (c *CategoryIndexer) Index(docs []Document) {
	for _, doc := range docs {
		cat := strings.ToLower(doc.Category)
		if _, exists := c.categoryIndex[cat]; !exists {
			c.categoryIndex[cat] = make(map[int]struct{})
		}
		c.categoryIndex[cat][doc.ID] = struct{}{}
	}
}

func (c *CategoryIndexer) GetDocsByCategories(categories []string) map[int]struct{} {
	result := make(map[int]struct{})
	for _, cat := range categories {
		for id := range c.categoryIndex[strings.ToLower(cat)] {
			result[id] = struct{}{}
		}
	}
	return result
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

func GetRankingStrategy(method string) RankingStrategy {
	switch method {
	case "size":
		return &ByDocSize{}
	case "frequency":
		return &ByKeywordFrequency{}
	default:
		return &ByDocSize{}
	}
}

// ====== Filter Strategy Pattern ======
type FilterStrategy interface {
	Filter(ids []int, docs map[int]Document) []int
}

type IndexedCategoryFilter struct {
	categoryIndexer *CategoryIndexer
	categories      []string
}

func NewIndexedCategoryFilter(indexer *CategoryIndexer, categories []string) *IndexedCategoryFilter {
	return &IndexedCategoryFilter{
		categoryIndexer: indexer,
		categories:      categories,
	}
}

func (f *IndexedCategoryFilter) Filter(ids []int, docs map[int]Document) []int {
	if len(f.categories) == 0 {
		return ids // no filtering
	}

	allowed := f.categoryIndexer.GetDocsByCategories(f.categories)
	filtered := make([]int, 0)
	for _, id := range ids {
		if _, exists := allowed[id]; exists {
			filtered = append(filtered, id)
		}
	}
	return filtered
}

// ====== Search Engine ======
type SearchEngine struct {
	documents       map[int]Document
	indexer         Indexer
	categoryIndexer *CategoryIndexer
}

func NewSearchEngine(indexer Indexer, catIndexer *CategoryIndexer) *SearchEngine {
	return &SearchEngine{
		documents:       make(map[int]Document),
		indexer:         indexer,
		categoryIndexer: catIndexer,
	}
}

func (s *SearchEngine) AddDocuments(docs []Document) {
	for _, doc := range docs {
		s.documents[doc.ID] = doc
	}
	s.indexer.Index(docs)
	s.categoryIndexer.Index(docs)
}

func (s *SearchEngine) Search(keyword, rankingMethod string, filter FilterStrategy) []Document {
	ids := s.indexer.Search(keyword)
	filtered := filter.Filter(ids, s.documents)
	ranker := GetRankingStrategy(rankingMethod)
	sortedIDs := ranker.Rank(filtered, s.documents, keyword)

	results := make([]Document, 0, len(sortedIDs))
	for _, id := range sortedIDs {
		results = append(results, s.documents[id])
	}
	return results
}

// ====== Main ======
func main() {
	docs := []Document{
		{ID: 1, Text: "Go is expressive, concise, clean, and efficient.", Category: "programming"},
		{ID: 2, Text: "Concurrency is not parallelism.", Category: "concepts"},
		{ID: 3, Text: "Go makes it easy to build simple, reliable, and efficient software.", Category: "programming"},
		{ID: 4, Text: "Software engineering is about trade-offs.", Category: "engineering"},
	}

	textIndexer := NewInvertedIndexer()
	categoryIndexer := NewCategoryIndexer()
	searchEngine := NewSearchEngine(textIndexer, categoryIndexer)
	searchEngine.AddDocuments(docs)

	filter := NewIndexedCategoryFilter(categoryIndexer, []string{"programming"})
	results := searchEngine.Search("efficient", "frequency", filter)

	fmt.Println("Search 'efficient' in category 'programming':")
	for _, doc := range results {
		fmt.Printf("Doc %d: %s (Category: %s)\n", doc.ID, doc.Text, doc.Category)
	}
}
