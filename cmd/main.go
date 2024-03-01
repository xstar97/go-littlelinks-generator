package main

import (
	"flag"
	"fmt"
	"github.com/xstar97/go-littlelinks-generator/internal/config"
	"github.com/xstar97/go-littlelinks-generator/internal/utils"
)

func main() {
	// Define flags for asset path and config file
	assetPathFlag := flag.String("asset-path", "assets/", "Path to the assets directory")
	configFlag := flag.String("config", "links.json", "Path to the configuration file")
	flag.Parse()

	assetsPath := *assetPathFlag
	linksPath := assetsPath + *configFlag

	// Parse the links JSON file
	links, err := utils.ParseLinksJSON(linksPath)
	if err != nil {
		fmt.Println("Error parsing links JSON:", err)
		return
	}
	utils.ParseConfig(links)

	// Delete the build directory if it exists
	deleteBuildDir()

	// Download the latest release
	err = utils.DownloadLatestRelease(links.DownloadTagVer)
	if err != nil {
		fmt.Println("Error downloading latest release:", err)
		return
	}

	// Generate HTML based on the configuration file
	generateHtml(links)

	// Check if redirects are enabled
	generateRedirect(links)

	// Copy any assets if found
	utils.ValidateAndCopyLinksAssets(links, assetsPath)

	// Delete files
	utils.CleanUpBuildFiles()
}

func deleteBuildDir() {
	if err := utils.DeleteBuildDirectory(config.BUILD_DIR); err != nil {
		fmt.Println("Error deleting build directory:", err)
		return
	}
}

func generateHtml(links *utils.Links) {
	// Generate HTML based on the configuration file
	genErr := utils.GenerateHTML(links)
	if genErr != nil {
		fmt.Println("Error generating HTML:", genErr)
		return
	}
	fmt.Println("\nHTML generation completed successfully.")
}

func generateRedirect(links *utils.Links) {
	// Check if redirects are enabled
	if links.EnableRedirects {
		fmt.Println("\nRedirects feature is enabled; generating _redirects...")
		// Generate redirects
		err := utils.GenerateRedirects(links)
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