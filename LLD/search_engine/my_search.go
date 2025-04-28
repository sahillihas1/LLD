package main

import "strings"

type NewSearchEnginee struct {
	Indexer           IIndexerr
	compositeSearcher ICompositeSearcher
	docs              map[int]Documentt
}

type Documentt struct {
	ID       int
	Text     string
	Category string
}

type ICompositeSearcher interface {
	Search(request []*searchRequest) []int
}

type searchRequest struct {
	key        string
	val        string
	searchType string
}

type CompositeSearcher struct {
	Searchers []ISearcher
}

func NewCompositeSearcher(searchers []ISearcher) *CompositeSearcher {
	return &CompositeSearcher{Searchers: searchers}
}

func (c *CompositeSearcher) Search(request []*searchRequest) []int {
	results := make(map[int]bool)
	for _, req := range request {
		for _, searcher := range c.Searchers {
			if searcher != nil {
				res := searcher.Search(req.key)
				for _, id := range res {
					results[id] = true
				}
			}
		}
	}

	var ids []int
	for id := range results {
		ids = append(ids, id)
	}
	return ids
}

type IIndexerr interface {
	Index(docs Documentt)
}

type IIndexerRepo interface {
	Index(docs Documentt)
	GetDocument(word string) Documentt
}

type IndexerRepo struct {
	index map[string][]int
}

func NewIndexerRepo() *IndexerRepo {
	return &IndexerRepo{index: make(map[string][]int)}
}

func (i *IndexerRepo) Index(docs Documentt) {
	words := strings.Fields(strings.ToLower(docs.Text))
	seen := make(map[string]bool)
	for _, word := range words {
		if !seen[word] {
			i.index[word] = append(i.index[word], docs.ID)
			seen[word] = true
		}
	}
}

func (i *IndexerRepo) GetDocument(word string) Documentt {
	if ids, exists := i.index[word]; exists {
		return Documentt{ID: ids[0], Text: word} // Simplified for example
	}
	return Documentt{}
}

type Indexerr struct {
	Indrepo IIndexerRepo
}

func (i *Indexerr) Index(docs Documentt) {
	i.Indrepo.Index(docs)
}

type ISearcher interface {
	Search(keyword string) []int
}

type FullTextSearcher struct {
	Indrepo IIndexerRepo
}

func (s *FullTextSearcher) Search(keyword string) []int {

}

type TermSearcher struct {
}

func (s *TermSearcher) Search(keyword string) []int {

}

func (i *Indexerr) Search(keyword string) []int {
	return i.index[keyword]
}

func main() {
	indexer := &Indexerr{Indrepo: NewIndexerRepo()}
	searchEngine := NewSearchEnginee{
		Indexer:           indexer,
		docs:              make(map[int]Documentt),
		compositeSearcher: NewCompositeSearcher([]ISearcher{&FullTextSearcher{}, &TermSearcher{}}),
	}

	docs := []Documentt{
		{ID: 1, Text: "Hello world", Category: "greeting"},
		{ID: 2, Text: "Goodbye world", Category: "farewell"},
	}

	for _, doc := range docs {
		indexer.Index(doc)
		searchEngine.docs[doc.ID] = doc
	}

	results := searchEngine.compositeSearcher.Search([]*searchRequest{
		{key: "Hello", val: "", searchType: "fulltext"},
		{key: "world", val: "", searchType: "term"},
	})
	for _, id := range results {
		doc := searchEngine.docs[id]
		println(doc.Text)
	}
}
