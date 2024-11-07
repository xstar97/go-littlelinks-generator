package utils

import (
	"fmt"
	//"io/ioutil"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/xstar97/go-littlelinks-generator/internal/config"
)

// GenerateHTML generates HTML from the template and JSON configuration
func GenerateHTML(conf *Config) error {

	//CopyEmbeddedFile() indexHTML
	buildPath := config.BUILD_DIR
	// Read the HTML template
	templateData, err := ReadIndexHTML()
	if err != nil {
		return err
	}

	// Replace template values with JSON matching keys
	htmlContent := string(templateData)
	htmlContent = strings.ReplaceAll(htmlContent, "{{META_ICON_URL}}", conf.Meta.IconURL)
	htmlContent = strings.ReplaceAll(htmlContent, "{{META_TITLE}}", conf.Meta.Title)
	htmlContent = strings.ReplaceAll(htmlContent, "{{META_AUTHOR}}", conf.Meta.Author)
	htmlContent = strings.ReplaceAll(htmlContent, "{{META_DESCRIPTION}}", conf.Meta.Description)
	htmlContent = strings.ReplaceAll(htmlContent, "{{META_THEME}}", conf.Meta.Theme)
	htmlContent = strings.ReplaceAll(htmlContent, "{{BIO_ICON_URL}}", conf.Bio.IconURL)
	htmlContent = strings.ReplaceAll(htmlContent, "{{BIO_TITLE}}", conf.Bio.Title)
	htmlContent = strings.ReplaceAll(htmlContent, "{{BIO_DESCRIPTION}}", conf.Bio.Description)
	htmlContent = strings.ReplaceAll(htmlContent, "{{BIO_FOOTER}}", conf.Bio.Footer)

	// Generate buttons for BIO_BUTTONS
	var bioButtons strings.Builder
	for _, link := range conf.Links {
		// Prepend BaseShortURL if it contains a URL
		if conf.BaseShortURL != "" {
			link.Link = conf.BaseShortURL + link.Link
		}
		fmt.Println("\nLink: ", link.Link)

		// Validate the button class
		if exists, err := ValidateButtonClass(link.Brand); err != nil || !exists {
			// Handle error or non-existence, fallback to default button class
			link.Brand = config.BUTTON_DETAILS_DEF_BRAND
		}
		fmt.Println("Brand: ", link.Brand)
		// Validate the button icon
		if exists, err := ValidateButtonImage(link.Icon); err != nil || !exists {
			// Handle error or non-existence, fallback to default button class
			link.Icon = config.BUTTON_DETAILS_DEF_ICON
		}
		fmt.Println("Icon: ", link.Icon)

		// Generate button HTML
		buttonHTML := fmt.Sprintf(`<a class="button button-%s" href="%s" target="_blank" rel="noopener" role="button"><img class="icon" src="images/icons/%s.svg" alt="">%s</a><br>`, link.Brand, link.Link, link.Icon, link.Name)
		// Add indentation before each button
		buttonHTML = "\t\t\t\t" + strings.ReplaceAll(buttonHTML, "\n", "\n\t")
		bioButtons.WriteString(buttonHTML)
		// Add a line break after each button
		bioButtons.WriteString("\n")
	}

	// Replace BIO_BUTTONS placeholder with generated buttons
	htmlContent = strings.ReplaceAll(htmlContent, "{{BIO_BUTTONS}}", bioButtons.String())

	// Create output directory if it doesn't exist
	err = os.MkdirAll(buildPath, os.ModePerm)
	if err != nil {
		return err
	}

	// Write the updated HTML content to the output file
	outputFile := filepath.Join(buildPath, "index.html")
	err = os.WriteFile(outputFile, []byte(htmlContent), 0644)
	if err != nil {
		return err
	}

	return nil
}

// validates if the given button class exists in the brands.css file
func ValidateButtonClass(name string) (bool, error) {
	path := fmt.Sprintf(config.BRANDS_CSS_FILE, config.BUILD_DIR)
	// Read the content of the brands.css file
	cssContent, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("failed to read brands.css file: %v", err)
	}

	// Check if the button class exists in the CSS content
	buttonClass := fmt.Sprintf(config.BUTTON_CLASS_NAME, name)
	return strings.Contains(string(cssContent), buttonClass), nil
}

func ValidateButtonImage(name string) (bool, error) {
	path := fmt.Sprintf(config.IMAGES_ICONS, config.BUILD_DIR, name)
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// verifies if the icon URLs in the links struct point to files and copies them to the build directory
func ValidateAndCopyLinksAssets(conf *Config, assetsPath string) error {
	if err := EnsureDirExists(config.BUILD_DIR); err != nil {
		return err
	}

	copyAsset := func(url string) error {
		if url == "" {
			return nil
		}

		srcPath := filepath.Join(assetsPath, url)
		destPath := filepath.Join(config.BUILD_DIR, url)

		// Open the source file
		srcFile, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("error opening source file: %v", err)
		}
		defer srcFile.Close()

		// Create the destination directory if it doesn't exist
		if err := EnsureDirExists(filepath.Dir(destPath)); err != nil {
			return err
		}

		// Create the destination file
		destFile, err := os.Create(destPath)
		if err != nil {
			return fmt.Errorf("error creating destination file: %v", err)
		}
		defer destFile.Close()

		// Copy the file contents
		if _, err := io.Copy(destFile, srcFile); err != nil {
			return fmt.Errorf("error copying file contents: %v", err)
		}

		return nil
	}

	if err := copyAsset(conf.Bio.IconURL); err != nil {
		return err
	}
	if err := copyAsset(conf.Meta.IconURL); err != nil {
		return err
	}

	return nil
}
