package cloudmail

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const defaultAPIBaseURL = "https://cloud.mail.ru/api/v2"

var (
	publicLinkIDRE = regexp.MustCompile(`/public/([^/?#]+/[^/?#]+)`)
	pageIDRE       = regexp.MustCompile(`pageId['"]*:\s*['"]*([^"'\\s,]+)`)
)

type folderAPIResponse struct {
	Body struct {
		Name string `json:"name"`
		List []struct {
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"list"`
	} `json:"body"`
}

type dispatcherAPIResponse struct {
	Body struct {
		WeblinkGet []struct {
			URL string `json:"url"`
		} `json:"weblink_get"`
	} `json:"body"`
}

// Config configures resolver network behavior.
type Config struct {
	Timeout   time.Duration
	Proxy     string
	ProxyAuth string
	UserAgent string
}

// Resolver resolves Cloud.Mail public links into direct file links.
type Resolver struct {
	client    *http.Client
	apiBase   string
	userAgent string
}

// NewResolver creates a new resolver instance.
func NewResolver(cfg Config) (*Resolver, error) {
	timeout := cfg.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Second
	}

	transport, ok := http.DefaultTransport.(*http.Transport)
	if !ok {
		return nil, errors.New("unexpected default transport type")
	}
	cloned := transport.Clone()

	if strings.TrimSpace(cfg.Proxy) != "" {
		proxyURL, err := buildProxyURL(cfg.Proxy, cfg.ProxyAuth)
		if err != nil {
			return nil, fmt.Errorf("invalid proxy config: %w", err)
		}
		cloned.Proxy = http.ProxyURL(proxyURL)
	}

	userAgent := strings.TrimSpace(cfg.UserAgent)
	if userAgent == "" {
		userAgent = "cmrd/0.1"
	}

	return &Resolver{
		client: &http.Client{
			Timeout:   timeout,
			Transport: cloned,
		},
		apiBase:   defaultAPIBaseURL,
		userAgent: userAgent,
	}, nil
}

func buildProxyURL(proxyValue string, proxyAuth string) (*url.URL, error) {
	proxyValue = strings.TrimSpace(proxyValue)
	if !strings.Contains(proxyValue, "://") {
		proxyValue = "http://" + proxyValue
	}
	parsed, err := url.Parse(proxyValue)
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(proxyAuth) != "" {
		parts := strings.SplitN(proxyAuth, ":", 2)
		if len(parts) == 2 {
			parsed.User = url.UserPassword(parts[0], parts[1])
		} else {
			parsed.User = url.User(parts[0])
		}
	}
	return parsed, nil
}

// Resolve resolves a list of public links.
func (r *Resolver) Resolve(ctx context.Context, links []string) ([]File, error) {
	var all []File
	for _, raw := range links {
		link := strings.TrimSpace(raw)
		if link == "" {
			continue
		}
		files, err := r.resolvePublicLink(ctx, link)
		if err != nil {
			return nil, fmt.Errorf("resolve %q: %w", link, err)
		}
		all = append(all, files...)
	}
	return all, nil
}

func (r *Resolver) resolvePublicLink(ctx context.Context, link string) ([]File, error) {
	linkID, err := parsePublicLinkID(link)
	if err != nil {
		return nil, err
	}

	pageID, err := r.getPageID(ctx, link)
	if err != nil {
		return nil, err
	}

	baseURL, err := r.getBaseURL(ctx, pageID)
	if err != nil {
		return nil, err
	}

	return r.walkFolder(ctx, linkID, "", pageID, baseURL)
}

func parsePublicLinkID(link string) (string, error) {
	matches := publicLinkIDRE.FindStringSubmatch(link)
	if len(matches) < 2 {
		return "", fmt.Errorf("wrong public link: %s", link)
	}
	return matches[1], nil
}

func (r *Resolver) getPageID(ctx context.Context, link string) (string, error) {
	body, err := r.doGet(ctx, link)
	if err != nil {
		return "", err
	}
	matches := pageIDRE.FindStringSubmatch(body)
	if len(matches) < 2 {
		return "", errors.New("page id not found")
	}
	return matches[1], nil
}

func (r *Resolver) getBaseURL(ctx context.Context, pageID string) (string, error) {
	values := url.Values{}
	values.Set("x-page-id", pageID)
	endpoint := fmt.Sprintf("%s/dispatcher?%s", r.apiBase, values.Encode())

	body, err := r.doGet(ctx, endpoint)
	if err != nil {
		return "", err
	}

	var response dispatcherAPIResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return "", fmt.Errorf("decode dispatcher response: %w", err)
	}
	if len(response.Body.WeblinkGet) == 0 || strings.TrimSpace(response.Body.WeblinkGet[0].URL) == "" {
		return "", errors.New("base URL not found")
	}
	return response.Body.WeblinkGet[0].URL, nil
}

func (r *Resolver) walkFolder(ctx context.Context, linkID string, parentFolder string, pageID string, baseURL string) ([]File, error) {
	values := url.Values{}
	values.Set("weblink", linkID)
	values.Set("x-page-id", pageID)
	endpoint := fmt.Sprintf("%s/folder?%s", r.apiBase, values.Encode())

	body, err := r.doGet(ctx, endpoint)
	if err != nil {
		return nil, err
	}

	var response folderAPIResponse
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		return nil, fmt.Errorf("decode folder response: %w", err)
	}

	currentFolder := joinPath(parentFolder, response.Body.Name)
	var files []File

	for _, item := range response.Body.List {
		switch item.Type {
		case "folder":
			childLink := joinPath(linkID, item.Name)
			childFiles, err := r.walkFolder(ctx, childLink, currentFolder, pageID, baseURL)
			if err != nil {
				return nil, err
			}
			files = append(files, childFiles...)
		default:
			outputPath := sanitizeWindowsPath(joinPath(currentFolder, item.Name))
			directURL := strings.TrimRight(baseURL, "/") + "/" + encodeURLPath(joinPath(linkID, item.Name))
			files = append(files, File{
				URL:    directURL,
				Output: outputPath,
			})
		}
	}

	return files, nil
}

func (r *Resolver) doGet(ctx context.Context, endpoint string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", r.userAgent)

	resp, err := r.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("http status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func joinPath(parts ...string) string {
	normalized := make([]string, 0, len(parts))
	for _, part := range parts {
		cleaned := strings.Trim(part, "/")
		if cleaned != "" {
			normalized = append(normalized, cleaned)
		}
	}
	return strings.Join(normalized, "/")
}

func encodeURLPath(path string) string {
	parts := strings.Split(path, "/")
	encoded := make([]string, 0, len(parts))
	for _, part := range parts {
		if part == "" {
			continue
		}
		encoded = append(encoded, url.PathEscape(part))
	}
	return strings.Join(encoded, "/")
}

func sanitizeWindowsPath(value string) string {
	var builder strings.Builder
	builder.Grow(len(value))
	for _, r := range value {
		if r >= 0 && r <= 31 {
			continue
		}
		switch r {
		case '<', '>', ':', '"', '|', '?', '*':
			continue
		default:
			builder.WriteRune(r)
		}
	}
	return builder.String()
}

// ParsePercent extracts percentage from aria2-like text output.
func ParsePercent(line string) (float64, bool) {
	for i := 0; i < len(line); i++ {
		if line[i] != '%' {
			continue
		}
		start := i - 1
		for start >= 0 && line[start] >= '0' && line[start] <= '9' {
			start--
		}
		start++
		if start >= i {
			continue
		}
		percent, err := strconv.ParseFloat(line[start:i], 64)
		if err == nil {
			return percent, true
		}
	}
	return 0, false
}
