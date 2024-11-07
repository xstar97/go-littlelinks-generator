package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Meta represents the meta section of the JSON/YAML
type Meta struct {
	Title       string `json:"TITLE" yaml:"TITLE"`
	Author      string `json:"AUTHOR" yaml:"AUTHOR"`
	Description string `json:"DESCRIPTION" yaml:"DESCRIPTION"`
	IconURL     string `json:"ICON_URL" yaml:"ICON_URL"`
	Theme       string `json:"THEME" yaml:"THEME"`
}

// Bio represents the bio section of the JSON/YAML
type Bio struct {
	IconURL     string `json:"ICON_URL" yaml:"ICON_URL"`
	Title       string `json:"TITLE" yaml:"TITLE"`
	Description string `json:"DESCRIPTION" yaml:"DESCRIPTION"`
	Footer      string `json:"FOOTER" yaml:"FOOTER"`
}

// Config represents the entire JSON/YAML structure
type Config struct {
	DownloadTagVer  string `json:"DOWNLOAD_TAG_VER" yaml:"DOWNLOAD_TAG_VER"`
	Meta            Meta   `json:"META" yaml:"META"`
	Bio             Bio    `json:"BIO" yaml:"BIO"`
	BaseShortURL    string `json:"BASE_SHORT_URL" yaml:"BASE_SHORT_URL"`
	EnableRedirects bool   `json:"ENABLE_REDIRECTS" yaml:"ENABLE_REDIRECTS"`
	Links           []Link `json:"LINKS" yaml:"LINKS"`
}

// Link represents a single link entry
type Link struct {
	Brand     string `json:"brand" yaml:"brand"`
	Icon      string `json:"icon" yaml:"icon"`
	Name      string `json:"name" yaml:"name"`
	Link      string `json:"link" yaml:"link"`
	Redirects []struct {
		Src  string `json:"src" yaml:"src"`
		Dest string `json:"dest" yaml:"dest"`
		Code int    `json:"code" yaml:"code"`
	} `json:"redirects" yaml:"redirects"`
}

// ParseConfigData parses the JSON or YAML file at the given path and returns the Config structure
func ParseConfigData(assetsPath string) (*Config, error) {
	data, err := os.ReadFile(assetsPath)
	if err != nil {
		return nil, err
	}

	var config Config
	switch filepath.Ext(assetsPath) {
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
	default:
		return nil, errors.New("unsupported file format")
	}

	if err != nil {
		return nil, err
	}

	return &config, nil
}
