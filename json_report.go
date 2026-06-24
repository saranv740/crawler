package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"github.com/saranv740/crawler/internal/crawler"
)

// writeJSONReport sorts the data based on keys and writes to the provided filePath
func writeJSONReport(pages map[string]crawler.PageData, pathToWrite string) error {
	keys := make([]string, 0, len(pages))
	for k := range pages {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	sorted := make([]crawler.PageData, 0, len(pages))
	for _, k := range keys {
		sorted = append(sorted, pages[k])
	}

	data, err := json.MarshalIndent(sorted, "", "	")
	if err != nil {
		return fmt.Errorf("error in marshaling the data")
	}

	cleanedPath := filepath.Clean(pathToWrite)
	err = os.WriteFile(cleanedPath, data, 0644)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("error in writing file")
	}

	return nil
}
