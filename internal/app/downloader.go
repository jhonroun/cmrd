package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cavaliergopher/grab/v3"
)

func downloadFile(fileinfo File) {
	// create client
	client := grab.NewClient()

	fileinfo.Link = strings.Replace(fileinfo.Link, ".\\", "", 1)
	fileinfo.Link = strings.Replace(fileinfo.Link, "\\", "/", -1)
	fileinfo.Link = strings.Replace(fileinfo.Link, "https:/", "https://", 1)
	req, _ := grab.NewRequest(filepath.Join(downloadToFolder, fileinfo.Out), fileinfo.Link)

	// start download
	fmt.Printf("Downloading %v...\n", req.URL())
	resp := client.Do(req)
	fmt.Printf("  %v\n", resp.HTTPResponse.Status)

	// start UI loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

Loop:
	for {
		select {
		case <-t.C:
			fmt.Printf("  transferred %v / %v bytes (%.2f%%)\n",
				resp.BytesComplete(),
				resp.Size(),
				100*resp.Progress())

		case <-resp.Done:
			// download is complete
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Verify downloaded file...")
	fmt.Printf("Size: %v bytes / %v, files eq: %v\n", resp.BytesComplete(), int(fileinfo.Size), resp.BytesComplete() == int64(fileinfo.Size))
	fmt.Printf("Download saved to ./%v \n", resp.Filename)

}
