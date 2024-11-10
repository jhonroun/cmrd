package app

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

type App struct {
	MRULinks []string
	Files    []Files
}

func (a *App) Run() {
	printAbout()
	fmt.Println("Retrive and prepare download links...")
	// Read links
	links, err := readLines(linksFile)
	if err != nil {
		fmt.Println("Error reading links file:", err)
		return
	}
	a.MRULinks = links

	for _, link := range a.MRULinks {
		link = strings.TrimSpace(link)
		if !strings.HasPrefix(link, "http") {
			continue
		}
		files, err := getFilesLink(link)
		if err != nil {
			fmt.Println(err)
			continue
		}

		a.Files = append(a.Files, Files{BaseURL: link, FileCount: len(files), Files: files})
	}

	if DirExists(downloadToFolder) {
		err := os.Mkdir(downloadToFolder, 0755)
		if err != nil {
			panic(err)
		}
	}

	for _, f := range a.Files {
		fmt.Printf("Start download from link:%s\n", f.BaseURL)
		fmt.Printf("Total files:%d\n", f.FileCount)
		for i, file := range f.Files {
			fmt.Printf("Downloading file %d of %d from link:%s\n", i+1, f.FileCount, f.BaseURL)
			downloadFile(file)
		}

	}
}

func getFilesLink(link string) ([]File, error) {
	linkID, pageID, baseURL, err := getLinkInfo(link)
	if err != nil {
		return nil, err
	}

	f := []File{}

	files, err := getAllFiles(linkID, baseURL, pageID, "")
	if err != nil {
		return nil, errors.New("can't find any files")
	}
	for _, fileInfo := range files {
		f = append(f, File{Out: fileInfo.Output, Directory: downloadToFolder, Link: fileInfo.Link, Size: fileInfo.Size, Hash: fileInfo.Hash})
	}
	return f, nil
}

// getLinkInfo retrieves link ID, page ID, and base URL for the download
func getLinkInfo(link string) (string, string, string, error) {
	re := regexp.MustCompile(`/public/([^/]+/[^/]+)`)
	matches := re.FindStringSubmatch(link)
	if len(matches) < 2 {
		return "", "", "", errors.New("invalid link format")
	}

	linkID := matches[1]
	pageID, err := getPageID(link)
	if err != nil {
		return "", "", "", err
	}

	baseURL, err := getBaseURL(pageID)
	if err != nil {
		return "", "", "", err
	}

	return linkID, pageID, baseURL, nil
}

// getAllFiles retrieves all files in a folder based on the link ID
func getAllFiles(linkID, baseURL, pageID, folder string) ([]CMFile, error) {
	url := fmt.Sprintf("%sfolder?weblink=%s&x-page-id=%s", mainURL, linkID, pageID)
	url = sanitizeURL(url)

	resp, err := makeRequest(url)
	if err != nil {
		return nil, err
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(resp, &jsonResponse); err != nil {
		return nil, err
	}

	body, ok := jsonResponse["body"].(map[string]interface{})
	if !ok || body["list"] == nil {
		return nil, nil
	}

	var files []CMFile
	list, _ := body["list"].([]interface{})

	for _, item := range list {
		itemMap, _ := item.(map[string]interface{})
		if itemMap["type"] == "folder" {
			subFiles, _ := getAllFiles(path.Join(linkID, itemMap["name"].(string)), baseURL, pageID, itemMap["name"].(string))
			files = append(files, subFiles...)
		} else {
			output := sanitizeFilePath(filepath.Join(folder, itemMap["name"].(string)))
			downloadURL := filepath.Join(baseURL, linkID, itemMap["name"].(string))
			files = append(files, CMFile{Link: sanitizeURL(downloadURL), Output: output, Size: itemMap["size"].(float64), Hash: itemMap["hash"].(string)})
		}
	}
	return files, nil
}

// sanitizeFilePath removes invalid characters for Windows paths
func sanitizeFilePath(filename string) string {
	illegalChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range illegalChars {
		filename = strings.ReplaceAll(filename, char, "")
	}
	return filename
}

// getPageID retrieves the page ID from the URL
func getPageID(url string) (string, error) {
	page, err := makeRequest(url)
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`pageId['"]*:['"]*([^"'\\s,]+)`)
	match := re.FindStringSubmatch(string(page))
	if len(match) > 1 {
		return match[1], nil
	}

	return "", errors.New("page ID not found")
}

// getBaseURL retrieves the base URL for the page ID
func getBaseURL(pageID string) (string, error) {
	url := fmt.Sprintf("%sdispatcher?x-page-id=%s", mainURL, pageID)
	resp, err := makeRequest(url)
	if err != nil {
		return "", err
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(resp, &jsonResponse); err != nil {
		return "", err
	}

	if body, ok := jsonResponse["body"].(map[string]interface{}); ok {
		if weblink, ok := body["weblink_get"].([]interface{}); ok && len(weblink) > 0 {
			url, _ := weblink[0].(map[string]interface{})["url"].(string)
			return url, nil
		}
	}

	return "", errors.New("can't get base URL")
}

// makeRequest makes an HTTP GET request to the provided URL
func makeRequest(url string) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)

	if proxy != "" {
		//proxyURL := "http://" + proxy
		//proxy := http.ProxyURL(proxyURL)
		//client.Transport = &http.Transport{Proxy: proxy}
		if proxyAuth != "" {
			auth := "Basic " + base64.StdEncoding.EncodeToString([]byte(proxyAuth))
			req.Header.Add("Proxy-Authorization", auth)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// readLines reads lines from a file and returns them as a slice
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
