package storage

import (
	"encoding/json"
	"os"

	"github.com/taqiyeddinedj/daar-project3/pkg/indexer"
)

// SaveToFile saves the index to a JSON file
func SaveToFile(idx *indexer.Indexer, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(idx)
}

// LoadFromFile loads the index from a JSON file
func LoadFromFile(filename string) (*indexer.Indexer, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var idx indexer.Indexer
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&idx)

	return &idx, err
}
