package main

import (
	"flag"
	"fmt"
	"path/filepath"

	"github.com/xstar97/go-littlelinks-generator/internal/config"
	"github.com/xstar97/go-littlelinks-generator/internal/utils"
)

func main() {
	// Define flags for asset path and config file
	assetPathFlag := flag.String("asset-path", "assets/", "Path to the assets directory")
	configFlag := flag.String("config", "links.json", "Path to the configuration file (JSON or YAML)")
	flag.Parse()

	assetsPath := *assetPathFlag
	linksPath := filepath.Join(assetsPath, *configFlag)

	// Parse the config file (supports JSON or YAML)
	conf, err := utils.ParseConfigData(linksPath)
	if err != nil {
		fmt.Printf("Error parsing config file (%s): %v\n", filepath.Ext(linksPath), err)
		return
	}

	// Process the config data
	utils.ParseConfig(conf)

	// Delete the build directory if it exists
	deleteBuildDir()

	// Download the latest release
	err = utils.DownloadLatestRelease(config.DOWNLOAD_TAG_DEF_VER)
	if err != nil {
		fmt.Println("Error downloading latest release:", err)
		return
	}

	// Generate HTML based on the configuration file
	generateHtml(conf)

	// Check if redirects are enabled
	generateRedirect(conf)

	// Copy any assets if found
	utils.ValidateAndCopyLinksAssets(conf, assetsPath)

	// Delete files
	utils.CleanUpBuildFiles()
}

func deleteBuildDir() {
	if err := utils.DeleteBuildDirectory(config.BUILD_DIR); err != nil {
		fmt.Println("Error deleting build directory:", err)
		return
	}
}

func generateHtml(conf *utils.Config) {
	// Generate HTML based on the configuration file
	genErr := utils.GenerateHTML(conf)
	if genErr != nil {
		fmt.Println("Error generating HTML:", genErr)
		return
	}
	fmt.Println("\nHTML generation completed successfully.")
}

func generateRedirect(conf *utils.Config) {
	// Check if redirects are enabled
	if conf.EnableRedirects {
		fmt.Println("\nRedirects feature is enabled; generating _redirects...")
		// Generate redirects
		err := utils.GenerateRedirects(conf)
		if err != nil {
			fmt.Println("Error generating redirects:", err)
			return
		}
		fmt.Println("\nRedirects file generated successfully.")
	} else {
		// Redirects feature is disabled
		fmt.Println("\nRedirects feature is disabled; skipping _redirects file generation.")
	}
}
