package client

import (
	"fmt"
	"net/url"

	"github.com/leechael/readwise-skills/internal/model"
)

// AuthCheck validates the API token. Returns nil on success.
func (c *Client) AuthCheck() error {
	resp, err := c.get("/api/v2/auth/")
	if err != nil {
		return err
	}
	return checkNoContent(resp)
}

// HighlightListParams are query parameters for listing highlights.
type HighlightListParams struct {
	PageSize        int
	Page            int
	BookID          int
	UpdatedLT       string
	UpdatedGT       string
	HighlightedAtLT string
	HighlightedAtGT string
}

func (p HighlightListParams) encode() string {
	v := url.Values{}
	if p.PageSize > 0 {
		v.Set("page_size", fmt.Sprintf("%d", p.PageSize))
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.BookID > 0 {
		v.Set("book_id", fmt.Sprintf("%d", p.BookID))
	}
	if p.UpdatedLT != "" {
		v.Set("updated__lt", p.UpdatedLT)
	}
	if p.UpdatedGT != "" {
		v.Set("updated__gt", p.UpdatedGT)
	}
	if p.HighlightedAtLT != "" {
		v.Set("highlighted_at__lt", p.HighlightedAtLT)
	}
	if p.HighlightedAtGT != "" {
		v.Set("highlighted_at__gt", p.HighlightedAtGT)
	}
	if encoded := v.Encode(); encoded != "" {
		return "?" + encoded
	}
	return ""
}

// HighlightList returns a paginated list of highlights.
func (c *Client) HighlightList(params HighlightListParams) (*model.PaginatedResponse[model.Highlight], error) {
	resp, err := c.get("/api/v2/highlights/" + params.encode())
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.PaginatedResponse[model.Highlight]](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// HighlightGet returns a single highlight by ID.
func (c *Client) HighlightGet(id int) (*model.Highlight, error) {
	resp, err := c.get(fmt.Sprintf("/api/v2/highlights/%d/", id))
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.Highlight](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// HighlightCreate creates one or more highlights.
func (c *Client) HighlightCreate(req model.HighlightCreateRequest) ([]model.HighlightCreateResponse, error) {
	resp, err := c.post("/api/v2/highlights/", req)
	if err != nil {
		return nil, err
	}
	return decodeResponse[[]model.HighlightCreateResponse](resp)
}

// HighlightUpdate updates a highlight by ID.
func (c *Client) HighlightUpdate(id int, req model.HighlightUpdateRequest) (*model.Highlight, error) {
	resp, err := c.patch(fmt.Sprintf("/api/v2/highlights/%d/", id), req)
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.Highlight](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// HighlightDelete deletes a highlight by ID.
func (c *Client) HighlightDelete(id int) error {
	resp, err := c.delete(fmt.Sprintf("/api/v2/highlights/%d/", id))
	if err != nil {
		return err
	}
	return checkNoContent(resp)
}

// HighlightTagList lists tags for a highlight.
func (c *Client) HighlightTagList(highlightID int) (*model.PaginatedResponse[model.Tag], error) {
	resp, err := c.get(fmt.Sprintf("/api/v2/highlights/%d/tags", highlightID))
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.PaginatedResponse[model.Tag]](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// HighlightTagAdd adds a tag to a highlight.
func (c *Client) HighlightTagAdd(highlightID int, name string) (*model.Tag, error) {
	resp, err := c.post(fmt.Sprintf("/api/v2/highlights/%d/tags/", highlightID), map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.Tag](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// HighlightTagUpdate updates a tag on a highlight.
func (c *Client) HighlightTagUpdate(highlightID, tagID int, name string) (*model.Tag, error) {
	resp, err := c.patch(fmt.Sprintf("/api/v2/highlights/%d/tags/%d", highlightID, tagID), map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.Tag](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// HighlightTagDelete deletes a tag from a highlight.
func (c *Client) HighlightTagDelete(highlightID, tagID int) error {
	resp, err := c.delete(fmt.Sprintf("/api/v2/highlights/%d/tags/%d", highlightID, tagID))
	if err != nil {
		return err
	}
	return checkNoContent(resp)
}

// BookListParams are query parameters for listing books.
type BookListParams struct {
	PageSize          int
	Page              int
	Category          string
	Source            string
	UpdatedLT         string
	UpdatedGT         string
	LastHighlightAtLT string
	LastHighlightAtGT string
}

func (p BookListParams) encode() string {
	v := url.Values{}
	if p.PageSize > 0 {
		v.Set("page_size", fmt.Sprintf("%d", p.PageSize))
	}
	if p.Page > 0 {
		v.Set("page", fmt.Sprintf("%d", p.Page))
	}
	if p.Category != "" {
		v.Set("category", p.Category)
	}
	if p.Source != "" {
		v.Set("source", p.Source)
	}
	if p.UpdatedLT != "" {
		v.Set("updated__lt", p.UpdatedLT)
	}
	if p.UpdatedGT != "" {
		v.Set("updated__gt", p.UpdatedGT)
	}
	if p.LastHighlightAtLT != "" {
		v.Set("last_highlight_at__lt", p.LastHighlightAtLT)
	}
	if p.LastHighlightAtGT != "" {
		v.Set("last_highlight_at__gt", p.LastHighlightAtGT)
	}
	if encoded := v.Encode(); encoded != "" {
		return "?" + encoded
	}
	return ""
}

// BookList returns a paginated list of books.
func (c *Client) BookList(params BookListParams) (*model.PaginatedResponse[model.Book], error) {
	resp, err := c.get("/api/v2/books/" + params.encode())
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.PaginatedResponse[model.Book]](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BookGet returns a single book by ID.
func (c *Client) BookGet(id int) (*model.Book, error) {
	resp, err := c.get(fmt.Sprintf("/api/v2/books/%d/", id))
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.Book](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BookTagList lists tags for a book.
func (c *Client) BookTagList(bookID int) (*model.PaginatedResponse[model.Tag], error) {
	resp, err := c.get(fmt.Sprintf("/api/v2/books/%d/tags", bookID))
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.PaginatedResponse[model.Tag]](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BookTagAdd adds a tag to a book.
func (c *Client) BookTagAdd(bookID int, name string) (*model.Tag, error) {
	resp, err := c.post(fmt.Sprintf("/api/v2/books/%d/tags/", bookID), map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.Tag](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BookTagUpdate updates a tag on a book.
func (c *Client) BookTagUpdate(bookID, tagID int, name string) (*model.Tag, error) {
	resp, err := c.patch(fmt.Sprintf("/api/v2/books/%d/tags/%d", bookID, tagID), map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.Tag](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// BookTagDelete deletes a tag from a book.
func (c *Client) BookTagDelete(bookID, tagID int) error {
	resp, err := c.delete(fmt.Sprintf("/api/v2/books/%d/tags/%d", bookID, tagID))
	if err != nil {
		return err
	}
	return checkNoContent(resp)
}

// ExportParams are query parameters for the export endpoint.
type ExportParams struct {
	UpdatedAfter   string
	IDs            string
	IncludeDeleted bool
	PageCursor     string
}

func (p ExportParams) encode() string {
	v := url.Values{}
	if p.UpdatedAfter != "" {
		v.Set("updatedAfter", p.UpdatedAfter)
	}
	if p.IDs != "" {
		v.Set("ids", p.IDs)
	}
	if p.IncludeDeleted {
		v.Set("includeDeleted", "true")
	}
	if p.PageCursor != "" {
		v.Set("pageCursor", p.PageCursor)
	}
	if encoded := v.Encode(); encoded != "" {
		return "?" + encoded
	}
	return ""
}

// Export returns highlights export with cursor pagination.
func (c *Client) Export(params ExportParams) (*model.ExportResponse, error) {
	resp, err := c.get("/api/v2/export/" + params.encode())
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.ExportResponse](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// DailyReview returns the daily review highlights.
func (c *Client) DailyReview() (*model.DailyReview, error) {
	resp, err := c.get("/api/v2/review/")
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.DailyReview](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
