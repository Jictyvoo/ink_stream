package outdirwriter

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
)

func MoveFirstFileToCoverFolder(directory string) error {
	var files []string

	// Walk through all files in the directory and its subdirectories
	err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})

	if err != nil {
		return err
	}

	// Sort files by name
	slices.Sort(files)

	// Check if there are files to process
	if len(files) == 0 {
		return fmt.Errorf("no files found in directory")
	}

	// Define the _0Cover folder path
	coverFolder := filepath.Join(directory, CoverDirSuffix)

	// Check if _0Cover folder is empty
	coverFiles, err := os.ReadDir(coverFolder)
	if err != nil {
		return err
	}

	if len(coverFiles) > 0 {
		// No post-processing done
		return nil
	}

	// Move the first file to the _0Cover folder
	firstFile := files[0]
	destination := filepath.Join(coverFolder, filepath.Base(firstFile))

	if err = os.Rename(firstFile, destination); err != nil {
		return err
	}

	fmt.Printf("Moved file: %s to %s\n", firstFile, destination)
	return nil
}
