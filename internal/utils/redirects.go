package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xstar97/go-littlelinks-generator/internal/config"
)

// GenerateRedirects generates a _redirects file based on the redirects in the links JSON
func GenerateRedirects(conf *Config) error {

	// Create the _redirects file
	redirectsFilePath := filepath.Join(config.BUILD_DIR, config.REDIRECTS_FILE)
	redirectsFile, err := os.Create(redirectsFilePath)
	if err != nil {
		return err
	}
	defer redirectsFile.Close()

	// Write the redirects to the _redirects file
	for _, link := range conf.Links {
		for _, redirect := range link.Redirects {
			fmt.Printf("%s %s %d\n", redirect.Src, redirect.Dest, redirect.Code)
			_, err := redirectsFile.WriteString(fmt.Sprintf("%s %s %d\n", redirect.Src, redirect.Dest, redirect.Code))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// WriteRedirectsToFile writes redirects to a _redirects file in the build directory
func WriteRedirectsToFile(redirects []string) error {
	// Create the build directory if it doesn't exist
	err := os.MkdirAll(config.BUILD_DIR, os.ModePerm)
	if err != nil {
		return err
	}

	// Create or open the _redirects file for writing
	filePath := filepath.Join(config.BUILD_DIR, config.REDIRECTS_FILE)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write each redirect string to the file
	for _, redirect := range redirects {
		_, err := file.WriteString(redirect + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
