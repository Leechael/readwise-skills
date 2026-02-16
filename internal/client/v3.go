package client

import (
	"fmt"
	"net/url"

	"github.com/leechael/readwise-skills/internal/model"
)

// ReaderListParams are query parameters for listing Reader documents.
type ReaderListParams struct {
	ID               string
	Category         string
	Location         string
	UpdatedAfter     string
	PageCursor       string
	Tags             []string
	Limit            int
	WithHtmlContent  bool
	WithRawSourceUrl bool
}

func (p ReaderListParams) encode() string {
	v := url.Values{}
	if p.ID != "" {
		v.Set("id", p.ID)
	}
	if p.Category != "" {
		v.Set("category", p.Category)
	}
	if p.Location != "" {
		v.Set("location", p.Location)
	}
	if p.UpdatedAfter != "" {
		v.Set("updatedAfter", p.UpdatedAfter)
	}
	if p.PageCursor != "" {
		v.Set("pageCursor", p.PageCursor)
	}
	for _, t := range p.Tags {
		v.Add("tag", t)
	}
	if p.Limit > 0 {
		v.Set("limit", fmt.Sprintf("%d", p.Limit))
	}
	if p.WithHtmlContent {
		v.Set("withHtmlContent", "true")
	}
	if p.WithRawSourceUrl {
		v.Set("withRawSourceUrl", "true")
	}
	if encoded := v.Encode(); encoded != "" {
		return "?" + encoded
	}
	return ""
}

// ReaderList returns a paginated list of Reader documents.
func (c *Client) ReaderList(params ReaderListParams) (*model.CursorPaginatedResponse[model.Document], error) {
	resp, err := c.get("/api/v3/list/" + params.encode())
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.CursorPaginatedResponse[model.Document]](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ReaderSave saves a document to Reader.
func (c *Client) ReaderSave(req model.DocumentSaveRequest) (*model.DocumentSaveResponse, error) {
	resp, err := c.post("/api/v3/save/", req)
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.DocumentSaveResponse](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ReaderUpdate updates a Reader document by ID.
func (c *Client) ReaderUpdate(id string, req model.DocumentUpdateRequest) error {
	resp, err := c.patch(fmt.Sprintf("/api/v3/update/%s/", id), req)
	if err != nil {
		return err
	}
	return checkNoContent(resp)
}

// ReaderDelete deletes a Reader document by ID.
func (c *Client) ReaderDelete(id string) error {
	resp, err := c.delete(fmt.Sprintf("/api/v3/delete/%s/", id))
	if err != nil {
		return err
	}
	return checkNoContent(resp)
}

// ReaderTagList lists all Reader tags.
func (c *Client) ReaderTagList() (*model.CursorPaginatedResponse[model.ReaderTag], error) {
	resp, err := c.get("/api/v3/tags/")
	if err != nil {
		return nil, err
	}
	result, err := decodeResponse[model.CursorPaginatedResponse[model.ReaderTag]](resp)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
