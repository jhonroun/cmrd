package cmrd

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"sync"

	"github.com/jhonroun/cmrd/internal/aria2"
	"github.com/jhonroun/cmrd/internal/cloudmail"
)

// Client provides library API for resolve/download workflows.
type Client struct {
	cfg      Config
	resolver *cloudmail.Resolver
	runner   *aria2.Runner
}

// New creates a new client.
func New(cfg Config) (*Client, error) {
	cfg = cfg.normalized()

	resolver, err := cloudmail.NewResolver(cloudmail.Config{
		Timeout:   cfg.HTTPTimeout,
		Proxy:     cfg.Proxy,
		ProxyAuth: cfg.ProxyAuth,
	})
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg:      cfg,
		resolver: resolver,
		runner:   aria2.NewRunner(cfg.Aria2Path),
	}, nil
}

// Config returns effective configuration.
func (c *Client) Config() Config {
	return c.cfg
}

// ReadLinksFile reads links from text file, skipping blank and comment lines.
func ReadLinksFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var links []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		links = append(links, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return links, nil
}

// Resolve resolves cloud links into a flat file list.
func (c *Client) Resolve(ctx context.Context, links []string) ([]FileTask, error) {
	files, err := c.resolver.Resolve(ctx, links)
	if err != nil {
		return nil, err
	}

	result := make([]FileTask, 0, len(files))
	for _, file := range files {
		result = append(result, FileTask{
			URL:    file.URL,
			Output: file.Output,
		})
	}
	return result, nil
}

// Download resolves links and runs aria2c.
func (c *Client) Download(ctx context.Context, links []string, onProgress ProgressHandler) error {
	files, err := c.Resolve(ctx, links)
	if err != nil {
		return err
	}

	if onProgress != nil {
		onProgress(ProgressEvent{
			Phase:          "resolve",
			Message:        "resolve complete",
			TotalFiles:     len(files),
			DoneFiles:      0,
			RemainingFiles: len(files),
			CurrentFile:    currentFileForIndex(files, 0),
		})
	}

	return c.DownloadResolved(ctx, files, onProgress)
}

// DownloadResolved runs aria2c for already resolved files.
func (c *Client) DownloadResolved(ctx context.Context, files []FileTask, onProgress ProgressHandler) error {
	if len(files) == 0 {
		return errors.New("empty file list")
	}

	temp, err := os.CreateTemp("", "cmrd-input-*.txt")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer temp.Close()
	if c.cfg.DeleteInputAfterDone {
		defer os.Remove(tempPath)
	}

	internalFiles := make([]cloudmail.File, 0, len(files))
	for _, file := range files {
		internalFiles = append(internalFiles, cloudmail.File{
			URL:    file.URL,
			Output: file.Output,
		})
	}

	if err := aria2.WriteInput(temp, internalFiles, c.cfg.DownloadDir); err != nil {
		return fmt.Errorf("write aria2 input: %w", err)
	}

	if err := temp.Close(); err != nil {
		return err
	}

	if onProgress != nil {
		onProgress(ProgressEvent{
			Phase:          "download",
			Message:        "download started",
			TotalFiles:     len(files),
			DoneFiles:      0,
			RemainingFiles: len(files),
			CurrentFile:    currentFileForIndex(files, 0),
		})
	}

	progress := newDownloadProgress(len(files))
	err = c.runner.Run(ctx, tempPath, c.cfg.Proxy, c.cfg.ProxyAuth, func(event aria2.ProgressEvent) {
		if onProgress == nil {
			return
		}
		doneFiles, remainingFiles, currentFile := progress.update(event, files)
		onProgress(ProgressEvent{
			Phase:          event.Phase,
			Percent:        event.Percent,
			Message:        event.Message,
			TotalFiles:     len(files),
			DoneFiles:      doneFiles,
			RemainingFiles: remainingFiles,
			CurrentFile:    currentFile,
			Done:           event.Done,
			Err:            event.Err,
		})
	})
	if err != nil {
		return err
	}

	if onProgress != nil {
		onProgress(ProgressEvent{
			Phase:          "download",
			Percent:        100,
			Message:        "download completed",
			TotalFiles:     len(files),
			DoneFiles:      len(files),
			RemainingFiles: 0,
			Done:           true,
		})
	}
	return nil
}

type downloadProgress struct {
	mu       sync.Mutex
	total    int
	doneSeen int
}

func newDownloadProgress(total int) *downloadProgress {
	return &downloadProgress{total: total}
}

func (p *downloadProgress) update(event aria2.ProgressEvent, files []FileTask) (int, int, string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if strings.Contains(strings.ToLower(event.Message), "download complete:") && p.doneSeen < p.total {
		p.doneSeen++
	}

	doneFiles := p.doneSeen
	if event.Percent > 0 {
		estimated := int(math.Floor(event.Percent * float64(p.total) / 100.0))
		if estimated > doneFiles {
			doneFiles = estimated
		}
	}
	if event.Done {
		doneFiles = p.total
	}
	if doneFiles < 0 {
		doneFiles = 0
	}
	if doneFiles > p.total {
		doneFiles = p.total
	}

	remainingFiles := p.total - doneFiles
	if remainingFiles < 0 {
		remainingFiles = 0
	}

	return doneFiles, remainingFiles, currentFileForIndex(files, doneFiles)
}

func currentFileForIndex(files []FileTask, index int) string {
	if index < 0 || index >= len(files) {
		return ""
	}
	return files[index].Output
}
