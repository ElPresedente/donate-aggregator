package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unsafe"

	"github.com/Masterminds/semver/v3"
	"golang.org/x/sys/windows"
)

func isValidVersion(version string) bool {
	pattern := `^\d+\.\d+\.\d+$`
	matched, _ := regexp.MatchString(pattern, version)
	return matched
}

func getFileVersion(exePath string) (string, error) {
	if runtime.GOOS != "windows" {
		return "", fmt.Errorf("version metadata reading is only supported on Windows")
	}

	// Get file info size
	var zeroHandle windows.Handle
	infoSize, err := windows.GetFileVersionInfoSize(exePath, &zeroHandle)
	if err != nil || infoSize == 0 {
		return "", fmt.Errorf("failed to get version info size: %v", err)
	}

	// Get version info data
	info := make([]byte, infoSize)
	if err := windows.GetFileVersionInfo(exePath, 0, infoSize, unsafe.Pointer(&info[0])); err != nil {
		return "", fmt.Errorf("failed to get version info: %v", err)
	}

	// Query version value
	var verInfo *windows.VS_FIXEDFILEINFO
	var verInfoLen uint32
	fixedInfoPtr := unsafe.Pointer(&verInfo)
	if err := windows.VerQueryValue(unsafe.Pointer(&info[0]), `\`, fixedInfoPtr, &verInfoLen); err != nil {
		return "", fmt.Errorf("failed to query version value: %v", err)
	}

	// Extract version numbers
	major := verInfo.FileVersionMS >> 16
	minor := verInfo.FileVersionMS & 0xFFFF
	patch := verInfo.FileVersionLS >> 16
	return fmt.Sprintf("%d.%d.%d", major, minor, patch), nil
}

func getLatestRelease() (string, string, string, error) {
	url := "https://api.github.com/repos/ElPresedente/donate-aggregator/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", "", fmt.Errorf("GitHub API request failed with status: %d", resp.StatusCode)
	}

	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", "", err
	}

	latestVersion := strings.TrimPrefix(release.TagName, "v")
	var downloadURL, assetName string
	for _, asset := range release.Assets {
		if strings.HasSuffix(strings.ToLower(asset.Name), ".exe") {
			downloadURL = asset.BrowserDownloadURL
			assetName = asset.Name
			break
		}
	}

	if downloadURL == "" || assetName == "" {
		return "", "", "", fmt.Errorf("no .exe file found in the latest release")
	}

	return latestVersion, downloadURL, assetName, nil
}

func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status %d", resp.StatusCode)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	var currentVersion string
	var err error

	// Get the directory of the current executable
	exePath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		os.Exit(1)
	}
	exeDir := filepath.Dir(exePath)

	// Try to get version from donate-aggregator.exe metadata
	assetName := "Donate.Aggregator.exe" // Default name, will be updated from release
	existingExePath := filepath.Join(exeDir, assetName)
	if _, err := os.Stat(existingExePath); err == nil {
		currentVersion, err = getFileVersion(existingExePath)
		if err != nil {
			fmt.Printf("Warning: Could not read version from %s: %v\n", assetName, err)
		}
	}

	// Fall back to command-line argument if metadata version is invalid or not found
	if currentVersion == "" || !isValidVersion(currentVersion) {
		if len(os.Args) != 2 {
			fmt.Println("Usage: version_checker.exe <version>")
			fmt.Println("Example: version_checker.exe 1.0.0")
			fmt.Println("Note: Version could not be read from donate-aggregator.exe metadata")
			os.Exit(1)
		}
		currentVersion = os.Args[1]
		if !isValidVersion(currentVersion) {
			fmt.Println("Invalid version format. Please use x.x.x (e.g., 1.0.0)")
			os.Exit(1)
		}
	}

	// Get latest release info
	latestVersion, downloadURL, assetName, err := getLatestRelease()
	if err != nil {
		fmt.Printf("Error fetching release data: %v\n", err)
		os.Exit(1)
	}

	// Update existingExePath with the actual asset name from the release
	existingExePath = filepath.Join(exeDir, assetName)

	currentSemVer, err := semver.NewVersion(currentVersion)
	if err != nil {
		fmt.Printf("Invalid current version format: %v\n", err)
		os.Exit(1)
	}

	latestSemVer, err := semver.NewVersion(latestVersion)
	if err != nil {
		fmt.Printf("Invalid repository version format: %v\n", err)
		os.Exit(1)
	}

	if latestSemVer.GreaterThan(currentSemVer) {
		fmt.Printf("New version %s is available! (Current: %s)\n", latestVersion, currentVersion)
		fmt.Println("Downloading from:", downloadURL)

		// Download to a temporary file
		tempFile := filepath.Join(exeDir, assetName+".new")
		err = downloadFile(downloadURL, tempFile)
		if err != nil {
			fmt.Printf("Error downloading new version: %v\n", err)
			os.Exit(1)
		}

		// Remove the existing .exe if it exists
		if _, err := os.Stat(existingExePath); err == nil {
			err = os.Remove(existingExePath)
			if err != nil {
				fmt.Printf("Error removing old executable: %v\n", err)
				os.Exit(1)
			}
		}

		// Rename the new .exe to the final name
		err = os.Rename(tempFile, existingExePath)
		if err != nil {
			fmt.Printf("Error replacing executable: %v\n", err)
			os.Exit(1)
		}

		// Set executable permissions
		err = os.Chmod(existingExePath, 0755)
		if err != nil {
			fmt.Printf("Error setting permissions: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully updated %s to version %s\n", assetName, latestVersion)
	} else {
		fmt.Printf("You are using the latest version: %s\n", currentVersion)
	}
}
