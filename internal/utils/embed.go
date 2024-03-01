package utils

import (
	"embed"
	"io"
	"github.com/xstar97/go-littlelinks-generator/internal/config"
)

//go:embed templates/index.html
var indexHTML embed.FS

func ReadIndexHTML() ([]byte, error) {
	// Open the embedded file
	file, err := indexHTML.Open(config.TEMPLATE_HTML)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the contents of the file
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}