package aria2

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/jhonroun/cmrd/internal/cloudmail"
)

var percentRE = regexp.MustCompile(`(\d{1,3})%`)

// ProgressEvent represents one aria2 progress update.
type ProgressEvent struct {
	Phase   string
	Percent float64
	Message string
	Done    bool
	Err     error
}

// Runner executes aria2c.
type Runner struct {
	BinaryPath string
}

// NewRunner creates a new aria2 runner.
func NewRunner(binaryPath string) *Runner {
	binaryPath = strings.TrimSpace(binaryPath)
	if binaryPath == "" {
		binaryPath = strings.TrimSpace(os.Getenv("CMRD_ARIA2C_PATH"))
	}
	if binaryPath == "" {
		binaryPath = "aria2c"
	}
	return &Runner{BinaryPath: binaryPath}
}

// WriteInput writes aria2 input file format for provided files.
func WriteInput(w io.Writer, files []cloudmail.File, downloadDir string) error {
	for _, file := range files {
		if _, err := fmt.Fprintf(w, "%s\n\tout=%s\n\tdir=%s\n", file.URL, file.Output, downloadDir); err != nil {
			return err
		}
	}
	return nil
}

// Run starts aria2c and forwards progress updates.
func (r *Runner) Run(ctx context.Context, inputFile string, proxy string, proxyAuth string, onUpdate func(ProgressEvent)) error {
	args := []string{
		"--file-allocation=none",
		"--max-connection-per-server=10",
		"--split=10",
		"--max-concurrent-downloads=10",
		"--summary-interval=1",
		"--continue=true",
		`--user-agent=Mozilla/5.0 (compatible; Firefox/3.6; Linux)`,
		"--input-file=" + inputFile,
	}

	proxy = strings.TrimSpace(proxy)
	if proxy != "" {
		proxyValue := proxy
		if strings.TrimSpace(proxyAuth) != "" {
			proxyValue = strings.TrimSpace(proxyAuth) + "@" + proxyValue
		}
		args = append(args, "--all-proxy="+proxyValue)
	}

	cmd := exec.CommandContext(ctx, r.BinaryPath, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if onUpdate != nil {
		onUpdate(ProgressEvent{Phase: "download", Message: "aria2c started"})
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		readOutput(stdout, onUpdate)
	}()
	go func() {
		defer wg.Done()
		readOutput(stderr, onUpdate)
	}()

	waitErr := cmd.Wait()
	wg.Wait()

	if waitErr != nil {
		if onUpdate != nil {
			onUpdate(ProgressEvent{
				Phase:   "download",
				Message: waitErr.Error(),
				Done:    true,
				Err:     waitErr,
			})
		}
		return waitErr
	}

	if onUpdate != nil {
		onUpdate(ProgressEvent{
			Phase:   "download",
			Percent: 100,
			Message: "aria2c finished",
			Done:    true,
		})
	}
	return nil
}

func readOutput(reader io.Reader, onUpdate func(ProgressEvent)) {
	if onUpdate == nil {
		io.Copy(io.Discard, reader)
		return
	}

	scanner := bufio.NewScanner(reader)
	scanner.Split(splitCRLF)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		event := ProgressEvent{
			Phase:   "download",
			Message: line,
		}
		if match := percentRE.FindStringSubmatch(line); len(match) == 2 {
			percent, err := strconv.ParseFloat(match[1], 64)
			if err == nil {
				event.Percent = percent
			}
		}
		onUpdate(event)
	}

	if err := scanner.Err(); err != nil {
		onUpdate(ProgressEvent{
			Phase:   "download",
			Message: err.Error(),
			Err:     err,
		})
	}
}

func splitCRLF(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for i, b := range data {
		if b == '\n' || b == '\r' {
			return i + 1, bytes.TrimRight(data[:i], "\r\n"), nil
		}
	}
	if atEOF && len(data) > 0 {
		return len(data), bytes.TrimRight(data, "\r\n"), nil
	}
	return 0, nil, nil
}
