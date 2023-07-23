package indexer

import "github.com/blevesearch/bleve/v2"


func PerformSearch(index bleve.Index, query string) ([]string, error) {
	searchRequest := bleve.NewSearchRequest(bleve.NewQueryStringQuery(query))
	searchResults, err := index.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	var titles []string
	for _, hit := range searchResults.Hits {
		docID := hit.ID
		titles = append(titles, docID)
	}

	return titles, nil
}