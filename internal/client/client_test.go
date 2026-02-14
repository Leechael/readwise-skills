package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/leechael/readwise-skills/internal/model"
)

func newTestServer(handler http.HandlerFunc) (*httptest.Server, *Client) {
	ts := httptest.NewServer(handler)
	c := New("test-token")
	c.SetBaseURL(ts.URL)
	return ts, c
}

func TestAuthCheck(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Token test-token" {
			t.Errorf("unexpected auth header: %s", r.Header.Get("Authorization"))
		}
		if r.URL.Path != "/api/v2/auth/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer ts.Close()

	if err := c.AuthCheck(); err != nil {
		t.Fatalf("AuthCheck failed: %v", err)
	}
}

func TestAuthCheckUnauthorized(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"detail":"Invalid token."}`))
	})
	defer ts.Close()

	err := c.AuthCheck()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T", err)
	}
	if apiErr.StatusCode != 401 {
		t.Errorf("expected 401, got %d", apiErr.StatusCode)
	}
}

func TestHighlightList(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/highlights/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("page_size") != "10" {
			t.Errorf("unexpected page_size: %s", r.URL.Query().Get("page_size"))
		}
		json.NewEncoder(w).Encode(model.PaginatedResponse[model.Highlight]{
			Count: 1,
			Results: []model.Highlight{
				{ID: 123, Text: "test highlight"},
			},
		})
	})
	defer ts.Close()

	result, err := c.HighlightList(HighlightListParams{PageSize: 10})
	if err != nil {
		t.Fatalf("HighlightList failed: %v", err)
	}
	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}
	if result.Results[0].ID != 123 {
		t.Errorf("expected ID 123, got %d", result.Results[0].ID)
	}
}

func TestHighlightGet(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/highlights/456/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(model.Highlight{ID: 456, Text: "a highlight"})
	})
	defer ts.Close()

	result, err := c.HighlightGet(456)
	if err != nil {
		t.Fatalf("HighlightGet failed: %v", err)
	}
	if result.ID != 456 {
		t.Errorf("expected ID 456, got %d", result.ID)
	}
}

func TestHighlightCreate(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var req model.HighlightCreateRequest
		json.NewDecoder(r.Body).Decode(&req)
		if len(req.Highlights) != 1 || req.Highlights[0].Text != "new highlight" {
			t.Errorf("unexpected request body")
		}
		json.NewEncoder(w).Encode([]model.HighlightCreateResponse{
			{ID: 789, Title: "Test Book", ModifiedHighlights: []int{100}},
		})
	})
	defer ts.Close()

	result, err := c.HighlightCreate(model.HighlightCreateRequest{
		Highlights: []model.HighlightCreateItem{
			{Text: "new highlight", Title: "Test Book"},
		},
	})
	if err != nil {
		t.Fatalf("HighlightCreate failed: %v", err)
	}
	if len(result) != 1 || result[0].ID != 789 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestHighlightDelete(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v2/highlights/123/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer ts.Close()

	if err := c.HighlightDelete(123); err != nil {
		t.Fatalf("HighlightDelete failed: %v", err)
	}
}

func TestBookList(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/books/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(model.PaginatedResponse[model.Book]{
			Count:   2,
			Results: []model.Book{{ID: 1, Title: "Book 1"}, {ID: 2, Title: "Book 2"}},
		})
	})
	defer ts.Close()

	result, err := c.BookList(BookListParams{})
	if err != nil {
		t.Fatalf("BookList failed: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("expected 2 books, got %d", result.Count)
	}
}

func TestBookGet(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/books/42/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(model.Book{ID: 42, Title: "The Book"})
	})
	defer ts.Close()

	result, err := c.BookGet(42)
	if err != nil {
		t.Fatalf("BookGet failed: %v", err)
	}
	if result.Title != "The Book" {
		t.Errorf("expected title 'The Book', got %s", result.Title)
	}
}

func TestExport(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/export/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(model.ExportResponse{
			Count: 1,
			Results: []model.ExportResult{
				{UserBookID: 1, Title: "Exported Book"},
			},
		})
	})
	defer ts.Close()

	result, err := c.Export(ExportParams{})
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if result.Count != 1 {
		t.Errorf("expected count 1, got %d", result.Count)
	}
}

func TestDailyReview(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/review/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(model.DailyReview{
			ReviewID:  1,
			ReviewURL: "https://readwise.io/review/1",
			Highlights: []model.ReviewHighlight{
				{ID: 1, Text: "review highlight"},
			},
		})
	})
	defer ts.Close()

	result, err := c.DailyReview()
	if err != nil {
		t.Fatalf("DailyReview failed: %v", err)
	}
	if result.ReviewID != 1 {
		t.Errorf("expected review ID 1, got %d", result.ReviewID)
	}
}

func TestReaderList(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/list/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(model.CursorPaginatedResponse[model.Document]{
			Count:   1,
			Results: []model.Document{{ID: "abc", Title: "Article"}},
		})
	})
	defer ts.Close()

	result, err := c.ReaderList(ReaderListParams{})
	if err != nil {
		t.Fatalf("ReaderList failed: %v", err)
	}
	if result.Count != 1 {
		t.Errorf("expected 1 document, got %d", result.Count)
	}
}

func TestReaderSave(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(model.DocumentSaveResponse{ID: "new-id", URL: "https://example.com"})
	})
	defer ts.Close()

	result, err := c.ReaderSave(model.DocumentSaveRequest{URL: "https://example.com"})
	if err != nil {
		t.Fatalf("ReaderSave failed: %v", err)
	}
	if result.ID != "new-id" {
		t.Errorf("expected ID 'new-id', got %s", result.ID)
	}
}

func TestReaderDelete(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v3/delete/doc-id/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer ts.Close()

	if err := c.ReaderDelete("doc-id"); err != nil {
		t.Fatalf("ReaderDelete failed: %v", err)
	}
}

func TestReaderTagList(t *testing.T) {
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v3/tags/" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(model.CursorPaginatedResponse[model.ReaderTag]{
			Count:   2,
			Results: []model.ReaderTag{{Key: "k1", Name: "tag1"}, {Key: "k2", Name: "tag2"}},
		})
	})
	defer ts.Close()

	result, err := c.ReaderTagList()
	if err != nil {
		t.Fatalf("ReaderTagList failed: %v", err)
	}
	if result.Count != 2 {
		t.Errorf("expected 2 tags, got %d", result.Count)
	}
}

func TestRateLimitRetry(t *testing.T) {
	attempts := 0
	ts, c := newTestServer(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts == 1 {
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
	defer ts.Close()

	if err := c.AuthCheck(); err != nil {
		t.Fatalf("expected retry to succeed: %v", err)
	}
	if attempts != 2 {
		t.Errorf("expected 2 attempts, got %d", attempts)
	}
}
