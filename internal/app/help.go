package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PrettyPrintFiles function for a slice of File structs
func PrettyPrintFiles(files []File) {
	if len(files) == 0 {
		fmt.Println("No files to display.")
		return
	}

	for i, file := range files {
		fmt.Printf("File #%d:\n%s\n", i+1, file.String())
		if i < len(files)-1 {
			fmt.Println(strings.Repeat("-", 40)) // Divider between files
		}
	}
}

func DirExists(dirname string) bool {
	// Get the current application's directory
	appDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return false
	}

	// Construct the full path to the target directory
	dirPath := filepath.Join(appDir, dirname)

	// Get information about the file/directory
	info, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		return false // Directory does not exist
	}

	// Check if the path is a directory
	return info.IsDir()
}

func sanitizeURL(url string) string {
	return strings.Replace(url, " ", "%20", -1)
}
