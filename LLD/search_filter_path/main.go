package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileFilter interface - ensures all filters implement the `Matches` method
type FileFilter interface {
	Matches(file os.FileInfo) bool
}

// NameFilter - filters files by name substring match
type NameFilter struct {
	Substring string
}

func (nf NameFilter) Matches(file os.FileInfo) bool {
	return strings.Contains(file.Name(), nf.Substring)
}

// SizeFilter - filters files by size constraint (greater than given size)
type SizeFilter struct {
	MinSize int64 // in bytes
}

func (sf SizeFilter) Matches(file os.FileInfo) bool {
	return file.Size() >= sf.MinSize
}

// CompositeFilter - applies multiple filters
type CompositeFilter struct {
	Filters []FileFilter
}

func (cf CompositeFilter) Matches(file os.FileInfo) bool {
	for _, filter := range cf.Filters {
		if !filter.Matches(file) {
			return false
		}
	}
	return true
}

// FileSearcher - handles searching files in a directory based on filters
type FileSearcher struct {
	RootDir string
	Filter  FileFilter
}

// Search - walks the directory and finds matching files
func (fs FileSearcher) Search() ([]string, error) {
	var matchedFiles []string

	err := filepath.Walk(fs.RootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && fs.Filter.Matches(info) {
			matchedFiles = append(matchedFiles, path)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return matchedFiles, nil
}

func main() {
	// Define search directory
	dir := "./testdir" // Replace with your directory path

	// Create filters
	nameFilter := NameFilter{Substring: ".txt"}
	sizeFilter := SizeFilter{MinSize: 1024} // Files greater than 1KB

	// Composite filter to combine multiple criteria
	compositeFilter := CompositeFilter{Filters: []FileFilter{nameFilter, sizeFilter}}

	// Create FileSearcher with filters
	searcher := FileSearcher{RootDir: dir, Filter: compositeFilter}

	// Perform search
	matchingFiles, err := searcher.Search()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Display results
	fmt.Println("Matching files:")
	for _, file := range matchingFiles {
		fmt.Println(file)
	}
}
