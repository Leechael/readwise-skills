package model

// CursorPaginatedResponse is the v3 cursor-based pagination wrapper.
type CursorPaginatedResponse[T any] struct {
	Count          int     `json:"count"`
	NextPageCursor *string `json:"nextPageCursor"`
	Results        []T     `json:"results"`
}

// Document represents a Reader document.
type Document struct {
	ID              string            `json:"id"`
	URL             string            `json:"url"`
	SourceURL       string            `json:"source_url"`
	Title           string            `json:"title"`
	Author          string            `json:"author"`
	Source          string            `json:"source"`
	Category        string            `json:"category"`
	Location        string            `json:"location"`
	Tags            map[string]string `json:"tags"`
	SiteName        string            `json:"site_name"`
	WordCount       int               `json:"word_count"`
	ReadingTime     string            `json:"reading_time"`
	CreatedAt       string            `json:"created_at"`
	UpdatedAt       string            `json:"updated_at"`
	Notes           string            `json:"notes"`
	PublishedDate   string            `json:"published_date"`
	Summary         string            `json:"summary"`
	ImageURL        string            `json:"image_url"`
	ParentID        *string           `json:"parent_id"`
	ReadingProgress float64           `json:"reading_progress"`
	FirstOpenedAt   *string           `json:"first_opened_at"`
	LastOpenedAt    *string           `json:"last_opened_at"`
	SavedAt         string            `json:"saved_at"`
	LastMovedAt     string            `json:"last_moved_at"`
}

// DocumentSaveRequest is the request body for saving a document to Reader.
type DocumentSaveRequest struct {
	URL           string   `json:"url"`
	HTML          string   `json:"html,omitempty"`
	ShouldClean   *bool    `json:"should_clean_html,omitempty"`
	Title         string   `json:"title,omitempty"`
	Author        string   `json:"author,omitempty"`
	Summary       string   `json:"summary,omitempty"`
	PublishedDate string   `json:"published_date,omitempty"`
	ImageURL      string   `json:"image_url,omitempty"`
	Location      string   `json:"location,omitempty"`
	Category      string   `json:"category,omitempty"`
	SavedUsing    string   `json:"saved_using,omitempty"`
	Tags          []string `json:"tags,omitempty"`
	Notes         string   `json:"notes,omitempty"`
}

// DocumentSaveResponse is the response from saving a document.
type DocumentSaveResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// DocumentUpdateRequest is the request body for updating a document.
type DocumentUpdateRequest struct {
	Title    *string `json:"title,omitempty"`
	Author   *string `json:"author,omitempty"`
	Summary  *string `json:"summary,omitempty"`
	Notes    *string `json:"notes,omitempty"`
	Location *string `json:"location,omitempty"`
	Category *string `json:"category,omitempty"`
}

// ReaderTag represents a Reader tag.
type ReaderTag struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}
