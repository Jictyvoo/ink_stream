package utils

import (
	"log"
	"os"
	"path/filepath"
	"slices"
)

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}

	return nil
}

func ListAllFiles(inputFolder string) []string {
	// Find all .cbz dirEntries in the input folder
	dirEntries, err := os.ReadDir(inputFolder)
	if err != nil {
		log.Fatalf("Failed to read input folder: %v", err)
	}

	var parentFolders = slices.Repeat([]string{inputFolder}, len(dirEntries))
	var filenameList = make([]string, 0, len(dirEntries))
	for i := 0; i < len(dirEntries); i++ {
		dirEntry := dirEntries[i]
		parentDir := parentFolders[i]

		fileAbsPath := filepath.Join(parentDir, dirEntry.Name())
		if dirEntry.IsDir() {
			subDirFiles, subErr := os.ReadDir(fileAbsPath)
			if subErr != nil {
				log.Fatalf("Failed to read sub dir: %v", subErr)
			}
			dirEntries = append(dirEntries, subDirFiles...)
			parentFolders = slices.Grow(parentFolders, len(subDirFiles))
			for range len(subDirFiles) {
				parentFolders = append(parentFolders, fileAbsPath)
			}
		} else {
			filenameList = append(filenameList, fileAbsPath)
		}
	}

	slices.Sort(filenameList)
	return filenameList
}
