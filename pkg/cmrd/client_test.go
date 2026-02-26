package cmrd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadLinksFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "links.txt")

	content := `# comment

https://cloud.mail.ru/public/9bFs/gVzxjU5uC
  https://cloud.mail.ru/public/3umo/mCi4k2ZTs
`
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	links, err := ReadLinksFile(filePath)
	if err != nil {
		t.Fatalf("ReadLinksFile returned error: %v", err)
	}

	if len(links) != 2 {
		t.Fatalf("unexpected links count: got=%d want=2", len(links))
	}
	if links[0] != "https://cloud.mail.ru/public/9bFs/gVzxjU5uC" {
		t.Fatalf("unexpected first link: %q", links[0])
	}
	if links[1] != "https://cloud.mail.ru/public/3umo/mCi4k2ZTs" {
		t.Fatalf("unexpected second link: %q", links[1])
	}
}
