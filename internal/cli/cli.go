package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/jhonroun/cmrd/internal/grpcapi"
	"github.com/jhonroun/cmrd/internal/tui"
	"github.com/jhonroun/cmrd/pkg/cmrd"
)

const Version = "1.0.0"

// Run executes CMRD CLI.
func Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		printRootHelp(os.Stdout)
		return nil
	}

	switch args[0] {
	case "help", "--help", "-h":
		printRootHelp(os.Stdout)
		return nil
	case "version", "--version", "-v":
		fmt.Println(Version)
		return nil
	case "resolve":
		return runResolve(ctx, args[1:])
	case "download":
		return runDownload(ctx, args[1:])
	case "serve-grpc":
		return runServeGRPC(ctx, args[1:])
	default:
		return fmt.Errorf("unknown command %q\n\n%s", args[0], rootHelpText)
	}
}

func runResolve(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("resolve", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	linksPath := fs.String("links", "links.txt", "Path to file with public links")
	jsonOutput := fs.Bool("json", false, "Print JSON output")
	timeout := fs.Duration("timeout", 30*time.Second, "HTTP timeout")
	proxy := fs.String("proxy", "", "Proxy host:port or URL")
	proxyAuth := fs.String("proxy-auth", "", "Proxy auth in user:pass format")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printResolveHelp(os.Stdout)
			return nil
		}
		return err
	}

	links, err := cmrd.ReadLinksFile(*linksPath)
	if err != nil {
		return err
	}

	cfg := cmrd.DefaultConfig()
	cfg.HTTPTimeout = *timeout
	cfg.Proxy = strings.TrimSpace(*proxy)
	cfg.ProxyAuth = strings.TrimSpace(*proxyAuth)

	client, err := cmrd.New(cfg)
	if err != nil {
		return err
	}

	files, err := client.Resolve(ctx, links)
	if err != nil {
		return err
	}

	if *jsonOutput {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(files)
	}

	fmt.Printf("Resolved files: %d\n", len(files))
	for _, file := range files {
		fmt.Printf("%s\n  out=%s\n\n", file.URL, file.Output)
	}
	return nil
}

func runDownload(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("download", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	linksPath := fs.String("links", "links.txt", "Path to file with public links")
	downloadDir := fs.String("dir", "downloads", "Download destination directory")
	aria2Path := fs.String("aria2c", "", "Path to aria2c binary (or CMRD_ARIA2C_PATH)")
	timeout := fs.Duration("timeout", 30*time.Second, "HTTP timeout")
	proxy := fs.String("proxy", "", "Proxy host:port or URL")
	proxyAuth := fs.String("proxy-auth", "", "Proxy auth in user:pass format")
	tuiMode := fs.Bool("tui", true, "Enable Bubble Tea TUI")
	keepInput := fs.Bool("keep-input", false, "Keep generated aria2 input file")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printDownloadHelp(os.Stdout)
			return nil
		}
		return err
	}

	links, err := cmrd.ReadLinksFile(*linksPath)
	if err != nil {
		return err
	}

	cfg := cmrd.DefaultConfig()
	cfg.Aria2Path = strings.TrimSpace(*aria2Path)
	cfg.DownloadDir = strings.TrimSpace(*downloadDir)
	cfg.HTTPTimeout = *timeout
	cfg.Proxy = strings.TrimSpace(*proxy)
	cfg.ProxyAuth = strings.TrimSpace(*proxyAuth)
	cfg.DeleteInputAfterDone = !*keepInput

	client, err := cmrd.New(cfg)
	if err != nil {
		return err
	}

	if *tuiMode {
		return tui.RunDownload(ctx, client, links)
	}

	return client.Download(ctx, links, func(event cmrd.ProgressEvent) {
		if event.Percent > 0 {
			fmt.Printf("[%s] %.1f%% %s\n", strings.ToUpper(event.Phase), event.Percent, event.Message)
			return
		}
		fmt.Printf("[%s] %s\n", strings.ToUpper(event.Phase), event.Message)
	})
}

func runServeGRPC(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("serve-grpc", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	address := fs.String("listen", ":50051", "gRPC listen address")
	downloadDir := fs.String("dir", "downloads", "Default download destination")
	aria2Path := fs.String("aria2c", "", "Path to aria2c binary (or CMRD_ARIA2C_PATH)")
	timeout := fs.Duration("timeout", 30*time.Second, "HTTP timeout")
	proxy := fs.String("proxy", "", "Proxy host:port or URL")
	proxyAuth := fs.String("proxy-auth", "", "Proxy auth in user:pass format")
	keepInput := fs.Bool("keep-input", false, "Keep generated aria2 input file")

	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			printServeGRPCHelp(os.Stdout)
			return nil
		}
		return err
	}

	cfg := cmrd.DefaultConfig()
	cfg.Aria2Path = strings.TrimSpace(*aria2Path)
	cfg.DownloadDir = strings.TrimSpace(*downloadDir)
	cfg.HTTPTimeout = *timeout
	cfg.Proxy = strings.TrimSpace(*proxy)
	cfg.ProxyAuth = strings.TrimSpace(*proxyAuth)
	cfg.DeleteInputAfterDone = !*keepInput

	service := grpcapi.NewServer(cfg)
	fmt.Printf("gRPC server listening on %s\n", *address)
	return grpcapi.Serve(ctx, *address, service)
}

func printRootHelp(w io.Writer) {
	fmt.Fprint(w, rootHelpText)
}

func printResolveHelp(w io.Writer) {
	fmt.Fprint(w, resolveHelpText)
}

func printDownloadHelp(w io.Writer) {
	fmt.Fprint(w, downloadHelpText)
}

func printServeGRPCHelp(w io.Writer) {
	fmt.Fprint(w, serveGRPCHelpText)
}

const rootHelpText = `CMRD - Cloud.Mail downloader rewritten in idiomatic Go

Usage:
  cmrd <command> [flags]

Commands:
  resolve      Resolve Cloud.Mail public links to direct file URLs
  download     Resolve links and start download with aria2c
  serve-grpc   Start gRPC API server (experimental; not fully tested)
  version      Print version
  help         Show this help

Examples:
  cmrd resolve --links links.txt
  cmrd download --links links.txt --dir downloads --tui=true
  cmrd serve-grpc --listen :50051

Environment:
  CMRD_ARIA2C_PATH   Path to aria2c binary (used when --aria2c is not set)
`

const resolveHelpText = `Usage:
  cmrd resolve [flags]

Flags:
  --links string       Path to links file (default "links.txt")
  --json               Print JSON output
  --timeout duration   HTTP timeout (default 30s)
  --proxy string       Proxy host:port or URL
  --proxy-auth string  Proxy auth in user:pass format
`

const downloadHelpText = `Usage:
  cmrd download [flags]

Flags:
  --links string       Path to links file (default "links.txt")
  --dir string         Download destination directory (default "downloads")
  --aria2c string      Path to aria2c binary (fallback: CMRD_ARIA2C_PATH or "aria2c")
  --timeout duration   HTTP timeout (default 30s)
  --proxy string       Proxy host:port or URL
  --proxy-auth string  Proxy auth in user:pass format
  --tui bool           Enable Bubble Tea TUI (default true)
  --keep-input         Keep generated aria2 input file
`

const serveGRPCHelpText = `Usage:
  cmrd serve-grpc [flags]

Note:
  gRPC mode is experimental and not fully tested yet.

Flags:
  --listen string      Listen address (default ":50051")
  --dir string         Default download destination directory (default "downloads")
  --aria2c string      Path to aria2c binary (fallback: CMRD_ARIA2C_PATH or "aria2c")
  --timeout duration   HTTP timeout (default 30s)
  --proxy string       Proxy host:port or URL
  --proxy-auth string  Proxy auth in user:pass format
  --keep-input         Keep generated aria2 input file
`
