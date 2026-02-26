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

func TestEncodeURLPath(t *testing.T) {
	got := encodeURLPath("M3Yv/83McATjtK/[swband.co] 1 раздел/file name.txt")
	want := "M3Yv/83McATjtK/%5Bswband.co%5D%201%20%D1%80%D0%B0%D0%B7%D0%B4%D0%B5%D0%BB/file%20name.txt"
	if got != want {
		t.Fatalf("encodeURLPath mismatch: got=%q want=%q", got, want)
	}
}

func TestSanitizeWindowsPath(t *testing.T) {
	got := sanitizeWindowsPath(`folder<bad>|name:?.txt`)
	want := "folderbadname.txt"
	if got != want {
		t.Fatalf("sanitizeWindowsPath mismatch: got=%q want=%q", got, want)
	}
}
