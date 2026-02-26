package cloudmail

import "testing"

func TestParsePublicLinkID(t *testing.T) {
	tests := []struct {
		name    string
		link    string
		wantID  string
		wantErr bool
	}{
		{
			name:   "valid public link",
			link:   "https://cloud.mail.ru/public/9bFs/gVzxjU5uC",
			wantID: "9bFs/gVzxjU5uC",
		},
		{
			name:    "invalid link",
			link:    "https://example.com/file",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := parsePublicLinkID(tc.link)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got none")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.wantID {
				t.Fatalf("unexpected id: got=%q want=%q", got, tc.wantID)
			}
		})
	}
}

func TestJoinPath(t *testing.T) {
	got := joinPath("a/", "/b", "c")
	want := "a/b/c"
	if got != want {
		t.Fatalf("joinPath mismatch: got=%q want=%q", got, want)
	}
}

func TestSanitizeWindowsPath(t *testing.T) {
	got := sanitizeWindowsPath(`folder<bad>|name:?.txt`)
	want := "folderbadname.txt"
	if got != want {
		t.Fatalf("sanitizeWindowsPath mismatch: got=%q want=%q", got, want)
	}
}
