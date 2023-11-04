package splitter

import (
	"encoding/json"
	"os"
)

// Split represents slice of package names that are supposed to be tested in the particular split
type Split []string

func StoreSplits(split []Split, path string) error {
	content, err := json.MarshalIndent(split, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

func LoadSplits(path string) ([]Split, error) {
	rawData, err := os.ReadFile(path)
	if err != nil {
		return []Split{}, err
	}
	var split []Split
	err = json.Unmarshal(rawData, &split)
	if err != nil {
		return []Split{}, err
	}
	return split, nil
}
