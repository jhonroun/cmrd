package app

import "fmt"

// CMFile represents a downloadable file with link and output path
type CMFile struct {
	Link   string
	Output string
	Size   float64
	Hash   string
}

type Files struct {
	BaseURL   string
	FileCount int
	Files     []File
}

type File struct {
	Out       string
	Directory string
	Link      string
	Size      float64
	Hash      string
}

// String method for human-readable output of a single File struct
func (f File) String() string {
	return fmt.Sprintf("File:\n  Out: %s\n  Directory: %s\n  Link: %s\n  Size: %f\n  Hash: %s\n", f.Out, f.Directory, f.Link, f.Size, f.Hash)
}
