package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/xstar97/go-littlelinks-generator/internal/config"
)

func DeleteBuildDirectory(outputPath string) error {
	// Check if the build directory exists
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		// Build directory doesn't exist, nothing to delete
		return nil
	}

	// Remove the build directory and its contents
	err := os.RemoveAll(outputPath)
	if err != nil {
		return err
	}

	fmt.Println("Build directory deleted successfully.")
	return nil
}

func ParseConfig(conf *Config) {
	// Print the parsed links as JSON string with 4-tab spaces
	linksJSON, err := json.MarshalIndent(conf, "", "    ") // 4 spaces for indentation
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println(string(linksJSON) + "\n")
}

func ReplaceFile(srcFile, destFile string) error {
	src, err := os.ReadFile(srcFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(destFile, src, 0644)
	if err != nil {
		return err
	}

	return nil
}

// deletes the files specified in FILES_TO_DELETE from the build directory
func CleanUpBuildFiles() error {
	for _, file := range config.FILES_TO_DELETE {
		filePath := filepath.Join(config.BUILD_DIR, file)
		err := os.Remove(filePath)
		if err != nil && !os.IsNotExist(err) {
			// If the error is not "file not found", return it
			return fmt.Errorf("error deleting file %s: %v", filePath, err)
		}
	}
	return nil
}

// checks if a directory exists and creates it if it doesn't.
func EnsureDirExists(dirPath string) error {
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}
