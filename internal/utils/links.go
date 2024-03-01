package utils

import (
	"encoding/json"
	"io/ioutil"
)

// Link represents a single link entry
type Link struct {
	Brand     string `json:"brand"`
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	Link      string `json:"link"`
	Redirects []struct {
		Src  string `json:"src"`
		Dest string `json:"dest"`
		Code int    `json:"code"`
	} `json:"redirects"`
}

// Links represents the entire JSON structure
type Links struct {
	DownloadTagVer string  `json:"DOWNLOAD_TAG_VER"`
	Meta           Meta    `json:"META"`
	Bio            Bio     `json:"BIO"`
	BaseShortURL   string  `json:"BASE_SHORT_URL"`
	EnableRedirects bool   `json:"ENABLE_REDIRECTS"`
	Links          []Link  `json:"LINKS"`
}

// Meta represents the meta section of the JSON
type Meta struct {
	Title       string `json:"TITLE"`
	Author      string `json:"AUTHOR"`
	Description string `json:"DESCRIPTION"`
	IconURL     string `json:"ICON_URL"`
	Theme       string `json:"THEME"`
}

// Bio represents the bio section of the JSON
type Bio struct {
	IconURL     string `json:"ICON_URL"`
	Title       string `json:"TITLE"`
	Description string `json:"DESCRIPTION"`
	Footer      string `json:"FOOTER"`
}

// ParseLinksJSON parses the JSON file at the given path and returns the Links structure
func ParseLinksJSON(assetsPath string) (*Links, error) {
	data, err := ioutil.ReadFile(assetsPath)
	if err != nil {
		return nil, err
	}

	var links Links
	err = json.Unmarshal(data, &links)
	if err != nil {
		return nil, err
	}

	return &links, nil
}
