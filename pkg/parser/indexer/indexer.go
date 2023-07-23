package indexer

import (
	"log"
	"sda-manager/pkg/db/models"
	"strconv"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

func GetIndex(verses []models.Verse, indexName string) (bleve.Index, error) {
	// Build index mapping
	mapping, err := BuildVerseMapping()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var index bleve.Index
	index, err = bleve.Open(indexName)
	if err != nil {
		log.Fatal(err)

		// Create new index
		index, err = bleve.New(indexName, mapping)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		
		// Index data
		err = IndexVerseData(index, verses)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
	}
	return index, nil
}

func IndexVerseData(index bleve.Index, verse []models.Verse) error {
	for _, v := range verse {
		err := index.Index(strconv.Itoa(int(v.ID)), v)
		if err != nil {
			return err
		}
	}
	return nil

}

func BuildVerseMapping() (mapping.IndexMapping, error) {
	verseMapping := bleve.NewDocumentMapping()

	textFieldMapping := bleve.NewTextFieldMapping()
	hymnIdFieldMapping := bleve.NewNumericFieldMapping()
	stanzaFieldMapping := bleve.NewNumericFieldMapping()
	isChorusFieldMappiing := bleve.NewBooleanFieldMapping()


	verseMapping.AddFieldMappingsAt("Text", textFieldMapping)
	verseMapping.AddFieldMappingsAt("HymnID", hymnIdFieldMapping)
	verseMapping.AddFieldMappingsAt("Stanza", stanzaFieldMapping)
	verseMapping.AddFieldMappingsAt("IsChorus", isChorusFieldMappiing)

	indexMapping := bleve.NewIndexMapping()
	indexMapping.AddDocumentMapping("Verse", verseMapping)

	return indexMapping, nil
}
