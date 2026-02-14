package model

import "time"

// PaginatedResponse is the standard v2 pagination wrapper.
type PaginatedResponse[T any] struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []T    `json:"results"`
}

// Tag represents a Readwise tag.
type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Highlight represents a Readwise highlight.
type Highlight struct {
	ID            int       `json:"id"`
	Text          string    `json:"text"`
	Note          string    `json:"note"`
	Location      int       `json:"location"`
	LocationType  string    `json:"location_type"`
	HighlightedAt *string   `json:"highlighted_at"`
	CreatedAt     time.Time `json:"created_at"`
	URL           *string   `json:"url"`
	Color         string    `json:"color"`
	Updated       time.Time `json:"updated"`
	BookID        int       `json:"book_id"`
	ExternalID    *string   `json:"external_id"`
	Tags          []Tag     `json:"tags"`
	EndLocation   *int      `json:"end_location"`
	ReadwiseURL   string    `json:"readwise_url"`
}

// HighlightCreateRequest is the request body for creating highlights.
type HighlightCreateRequest struct {
	Highlights []HighlightCreateItem `json:"highlights"`
}

// HighlightCreateItem is a single highlight to create.
type HighlightCreateItem struct {
	Text         string  `json:"text"`
	Title        string  `json:"title,omitempty"`
	Author       string  `json:"author,omitempty"`
	ImageURL     string  `json:"image_url,omitempty"`
	SourceURL    string  `json:"source_url,omitempty"`
	SourceType   string  `json:"source_type,omitempty"`
	Category     string  `json:"category,omitempty"`
	Note         string  `json:"note,omitempty"`
	Location     *int    `json:"location,omitempty"`
	LocationType string  `json:"location_type,omitempty"`
	HighlightedAt string `json:"highlighted_at,omitempty"`
	HighlightURL string  `json:"highlight_url,omitempty"`
}

// HighlightUpdateRequest is the request body for updating a highlight.
type HighlightUpdateRequest struct {
	Text     *string `json:"text,omitempty"`
	Note     *string `json:"note,omitempty"`
	Location *string `json:"location,omitempty"`
	URL      *string `json:"url,omitempty"`
	Color    *string `json:"color,omitempty"`
}

// Book represents a Readwise book/source.
type Book struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Category        string    `json:"category"`
	Source          string    `json:"source"`
	NumHighlights   int       `json:"num_highlights"`
	LastHighlightAt *string   `json:"last_highlight_at"`
	Updated         time.Time `json:"updated"`
	CoverImageURL   string    `json:"cover_image_url"`
	HighlightsURL   string    `json:"highlights_url"`
	SourceURL       *string   `json:"source_url"`
	ASIN            *string   `json:"asin"`
	Tags            []Tag     `json:"tags"`
	DocumentNote    string    `json:"document_note"`
}

// HighlightCreateResponse is the response from creating highlights (returns book-like objects).
type HighlightCreateResponse struct {
	ID                 int       `json:"id"`
	Title              string    `json:"title"`
	Author             string    `json:"author"`
	Category           string    `json:"category"`
	Source             string    `json:"source"`
	NumHighlights      int       `json:"num_highlights"`
	LastHighlightAt    *string   `json:"last_highlight_at"`
	Updated            time.Time `json:"updated"`
	CoverImageURL      string    `json:"cover_image_url"`
	HighlightsURL      string    `json:"highlights_url"`
	SourceURL          *string   `json:"source_url"`
	ASIN               *string   `json:"asin"`
	Tags               []Tag     `json:"tags"`
	DocumentNote       string    `json:"document_note"`
	ModifiedHighlights []int     `json:"modified_highlights"`
}

// ExportResponse is the response from the export endpoint (cursor-paginated).
type ExportResponse struct {
	Count          int            `json:"count"`
	NextPageCursor *string        `json:"nextPageCursor"`
	Results        []ExportResult `json:"results"`
}

// ExportResult is a single book with its highlights from the export endpoint.
type ExportResult struct {
	UserBookID    int              `json:"user_book_id"`
	IsDeleted     bool             `json:"is_deleted"`
	Title         string           `json:"title"`
	Author        string           `json:"author"`
	ReadableTitle string           `json:"readable_title"`
	Source        string           `json:"source"`
	CoverImageURL string           `json:"cover_image_url"`
	UniqueURL     string           `json:"unique_url"`
	BookTags      []Tag            `json:"book_tags"`
	Category      string           `json:"category"`
	DocumentNote  string           `json:"document_note"`
	Summary       string           `json:"summary"`
	ReadwiseURL   string           `json:"readwise_url"`
	SourceURL     string           `json:"source_url"`
	ExternalID    *string          `json:"external_id"`
	ASIN          *string          `json:"asin"`
	Highlights    []ExportHighlight `json:"highlights"`
}

// ExportHighlight is a highlight within an export result.
type ExportHighlight struct {
	ID            int     `json:"id"`
	IsDeleted     bool    `json:"is_deleted"`
	Text          string  `json:"text"`
	Location      int     `json:"location"`
	LocationType  string  `json:"location_type"`
	Note          *string `json:"note"`
	Color         string  `json:"color"`
	HighlightedAt string  `json:"highlighted_at"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
	ExternalID    string  `json:"external_id"`
	EndLocation   *int    `json:"end_location"`
	URL           *string `json:"url"`
	BookID        int     `json:"book_id"`
	Tags          []Tag   `json:"tags"`
	IsFavorite    bool    `json:"is_favorite"`
	IsDiscard     bool    `json:"is_discard"`
	ReadwiseURL   string  `json:"readwise_url"`
}

// DailyReview is the response from the daily review endpoint.
type DailyReview struct {
	ReviewID        int               `json:"review_id"`
	ReviewURL       string            `json:"review_url"`
	ReviewCompleted bool              `json:"review_completed"`
	Highlights      []ReviewHighlight `json:"highlights"`
}

// ReviewHighlight is a highlight within a daily review.
type ReviewHighlight struct {
	ID            int     `json:"id"`
	Text          string  `json:"text"`
	Title         string  `json:"title"`
	Author        string  `json:"author"`
	Note          string  `json:"note"`
	Location      int     `json:"location"`
	LocationType  string  `json:"location_type"`
	Category      *string `json:"category"`
	SourceType    string  `json:"source_type"`
	HighlightedAt string  `json:"highlighted_at"`
	ImageURL      string  `json:"image_url"`
	HighlightURL  *string `json:"highlight_url"`
	URL           *string `json:"url"`
	SourceURL     *string `json:"source_url"`
}
