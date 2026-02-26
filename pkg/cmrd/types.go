package cmrd

// FileTask represents one download target.
type FileTask struct {
	URL    string `json:"url"`
	Output string `json:"output"`
}

// ProgressEvent is emitted during resolve/download lifecycle.
type ProgressEvent struct {
	Phase      string  `json:"phase"`
	Percent    float64 `json:"percent"`
	Message    string  `json:"message"`
	TotalFiles int     `json:"total_files"`
	Done       bool    `json:"done"`
	Err        error   `json:"-"`
}

// ProgressHandler receives progress events.
type ProgressHandler func(ProgressEvent)
