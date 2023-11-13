package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

func setGoEnvVariables(homeV, shellV string) error {
	shellVV := "." + shellV
	zshrcPath := filepath.Join(os.Getenv(homeV), shellVV)
	file, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening %s %w", shellV, err)
	}
	defer file.Close()

	_, err = file.WriteString("\nexport PATH=$PATH:/usr/local/go/bin\n")
	if err != nil {
		return fmt.Errorf("error writing to %s: %w", shellV, err)
	}
	_, err = file.WriteString("export GOROOT=/usr/local/go\n")
	if err != nil {
		return fmt.Errorf("error writing to %s: %w", shellV, err)
	}
	_, err = file.WriteString("export GOPATH=$HOME/go\n")
	if err != nil {
		return fmt.Errorf("error writing to %s: %w", shellV, err)
	}

	return nil
}

func checkHome() string {
	// Get the SHELL environment variable
	home := os.Getenv("HOME")

	// Check if the SHELL variable is set
	if home == "" {
		return "Unknown shell"
	}

	return home
}

// checkShell retrieves the SHELL environment variable and returns the shell name
func checkShell() string {
	// Get the SHELL environment variable
	fullShellPath := os.Getenv("SHELL")

	// Use filepath.Base to get the last element of the path
	shell := filepath.Base(fullShellPath)

	// Check if the SHELL variable is set
	if shell == "" {
		return "Unknown shell"
	}

	return shell
}

// deleteFile deletes a file at the given path
func deleteFile(filePath string) error {
	// Delete the file using os.Remove
	err := os.Remove(filePath)
	if err != nil {
		return err
	}
	return nil
}

func lastGoVer() (string, error) {
	// URL of the Go downloads page
	url := "https://go.dev/dl/"

	// Make a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching URL: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	// Regular expression to find Go version (only the numeric part)
	re := regexp.MustCompile(`go([0-9]+\.[0-9]+(\.[0-9]+)?)`)
	matches := re.FindStringSubmatch(string(body))

	if len(matches) < 2 {
		return "", fmt.Errorf("no Go versions found")
	}

	// The first submatch is the numeric part of the version
	latestVersion := matches[1]
	return latestVersion, nil
}

// extractTarGz extracts a .tar.gz file to a destination directory.
func extractTarGz(gzipStream io.Reader, dest string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()

	tarReader := tar.NewReader(uncompressedStream)
	for {
		header, err := tarReader.Next()

		switch {
		// If no more files are found return
		case err == io.EOF:
			return nil

		// Return any other error
		case err != nil:
			return err

		// If the header is nil, skip it
		case header == nil:
			continue
		}

		// The target location where the dir/file should be created
		target := filepath.Join(dest, header.Name)

		// Check the file type
		switch header.Typeflag {
		// If it's a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// If it's a file create it
		case tar.TypeReg:
			outFile, err := os.Create(target)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()
		}
	}
}

func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func main() {
	// Check if Go is installed and get version
	goVersion, err := exec.Command("go", "version").Output()
	if err != nil {
		fmt.Println("Go is not installed")
	} else {
		version := strings.Fields(string(goVersion))[2]
		fmt.Println("Go version:", version)
	}

	// Get CPU architecture
	cpuArch := runtime.GOARCH
	fmt.Println("CPU Architecture:", cpuArch)

	// Get Operating System
	osType := runtime.GOOS
	fmt.Println("Operating System:", osType)

	// Get latest Go version using Linux curl command

	latestGoVersion, err := lastGoVer()
	fmt.Println(latestGoVersion)

	if err != nil {
		fmt.Println("Error fetching latest Go version:", err)
		return
	}
	fmt.Println("Latest Go version:", string(latestGoVersion))

	// Construct the download URL
	downloadURL := fmt.Sprintf("https://dl.google.com/go/go%s.%s-%s.tar.gz", latestGoVersion, osType, cpuArch)
	filepath := fmt.Sprintf("go%s.%s-%s.tar.gz", latestGoVersion, osType, cpuArch)
	fmt.Println("Downloading file...")
	if err := downloadFile(downloadURL, filepath); err != nil {
		panic(err)
	}

	fmt.Println("Download completed successfully!")

	tarFilePath := filepath
	destPath := "/usr/local/"

	newFile, err := os.Open(tarFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer newFile.Close()

	if err := extractTarGz(newFile, destPath); err != nil {
		fmt.Println("Error extracting file:", err)
		return
	}

	fmt.Println("Extraction completed successfully!")

	// Call deleteFile function
	err = deleteFile(tarFilePath)
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}

	fmt.Println("The Go tar file successfully deleted")

	shell := checkShell()
	fmt.Println("Current shell:", shell)

	h := checkHome()
	fmt.Println(h)

	if err := setGoEnvVariables(h, shell); err != nil {
		fmt.Println("Error setting Go environment variables:", err)
		return
	}

}
