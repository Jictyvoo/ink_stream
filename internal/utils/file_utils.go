package utils

import (
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
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

	parentFolders := slices.Repeat([]string{inputFolder}, len(dirEntries))
	filenameList := make([]string, 0, len(dirEntries))
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

// CollapseFilesByExt takes a list of file paths and a list of file
// extensions. It returns a new list where any folder that contains *only*
// files with the supplied extensions is represented by the folder path
// itself. All other files are returned unchanged.
func CollapseFilesByExt(fileList []string, extList []string) []string {
	for index, fileExt := range extList {
		fileExt = strings.ToLower(strings.TrimSpace(fileExt))
		extList[index] = "." + strings.Trim(fileExt, ".") // ensure that it starts with a dot
	}

	// Group files by parent directory
	dirMap := make(map[string][]string)
	for _, elemFile := range fileList {
		dir := filepath.Dir(elemFile)
		dirMap[dir] = append(dirMap[dir], elemFile)
	}

	result := fileList[:0] // Reuse underlying array to reduce memory allocations
	for dir, files := range dirMap {
		allAllowed := true
		for _, subFile := range files {
			ext := strings.ToLower(filepath.Ext(subFile))
			if !slices.Contains(extList, ext) {
				allAllowed = false
				break
			}
		}
		if allAllowed {
			result = append(result, dir)
			continue
		}
		result = append(result, files...)
	}

	// Final processing after all valid folders identified
	slices.Sort(result)
	fileList = make([]string, 0, len(result))
	for index, filElem := range slices.Backward(result) {
		subSlice := result[:index]
		ext := filepath.Ext(filElem)
		folderDir := filepath.Dir(filElem)
		if ext != "" || !slices.Contains(subSlice, folderDir) {
			fileList = append(fileList, filElem)
		}
	}

	return fileList
}
