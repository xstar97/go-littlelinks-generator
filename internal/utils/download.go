package utils

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"github.com/google/go-github/github"
	"github.com/xstar97/go-littlelinks-generator/internal/config"
)

const (
	githubOwner = config.GH_REPO_OWNER
	githubRepo  = config.GH_REPO_NAME
)

func DownloadLatestRelease(versionTag string) error {
	ctx := context.Background()
	client := github.NewClient(nil)
	tempPath := config.TEMP_DIR
	buildPath := config.BUILD_DIR

	if versionTag == "" {
		release, _, err := client.Repositories.GetLatestRelease(ctx, githubOwner, githubRepo)
		if err != nil {
			return err
		}
		versionTag = *release.TagName
	}

	// Ensure temp and build directories exist
	if err := EnsureDirExists(tempPath); err != nil {
		return err
	}
	if err := EnsureDirExists(buildPath); err != nil {
		return err
	}

	zipFileName := fmt.Sprintf(config.DOWNLOAD_ZIP_NAME, versionTag)
	tempZipFile := filepath.Join(tempPath, zipFileName)
	srcDir := filepath.Join(tempPath, fmt.Sprintf("littlelink-%s", strings.TrimPrefix(versionTag, "v")))

	// Check if the zip file already exists
	if _, err := os.Stat(tempZipFile); err == nil {
		fmt.Printf("The zip file %s already exists. Skipping download.\n", zipFileName)
	} else {
		// Download the zip file
		if err := downloadZip(versionTag, tempZipFile); err != nil {
			return err
		}
	}

	// Check if the extracted directory already exists
	if _, err := os.Stat(srcDir); err == nil {
		fmt.Printf("The directory %s already exists. Skipping extraction.\n", srcDir)
	} else {
		// Extract the zip file into the temp directory
		if err := unzip(tempZipFile, tempPath); err != nil {
			return err
		}
	}

	// Copy the files and directories from the extracted directory to the build directory
	if err := copyDir(srcDir, buildPath); err != nil {
		return err
	}

	// Log extraction success
	fmt.Printf("Copied files from %s to %s\n", srcDir, buildPath)

	return nil
}

func downloadZip(versionTag, tempZipFile string) error {
	url := fmt.Sprintf("https://github.com/%s/%s/archive/refs/tags/%s.zip", githubOwner, githubRepo, versionTag)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download: %s", resp.Status)
	}

	out, err := os.Create(tempZipFile)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Log download success
	fmt.Printf("Downloaded %s into %s\n", tempZipFile, versionTag)

	return nil
}

// unzip extracts a zip file into the specified directory
func unzip(src string, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		if err = os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		dstFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(dstFile, rc)
		rc.Close()
		dstFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// copyDir copies files and directories from source to destination
func copyDir(src string, dest string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			return err
		}

		return nil
	})
}
