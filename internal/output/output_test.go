package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type testItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func TestJSONOutput(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(true, false, "")
	f.Writer = &buf

	data := testItem{ID: 1, Name: "test"}
	if err := f.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	var result testItem
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if result.ID != 1 || result.Name != "test" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestJQFilter(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(true, false, ".name")
	f.Writer = &buf

	data := testItem{ID: 1, Name: "hello"}
	if err := f.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := strings.TrimSpace(buf.String())
	if output != `"hello"` {
		t.Errorf("expected '\"hello\"', got %s", output)
	}
}

func TestJQFilterArray(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(true, false, ".[].name")
	f.Writer = &buf

	data := []testItem{{ID: 1, Name: "a"}, {ID: 2, Name: "b"}}
	if err := f.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d: %v", len(lines), lines)
	}
}

func TestPlainStructOutput(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(false, true, "")
	f.Writer = &buf

	data := testItem{ID: 42, Name: "myname"}
	if err := f.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "id") || !strings.Contains(output, "42") {
		t.Errorf("expected key-value output, got: %s", output)
	}
	if !strings.Contains(output, "name") || !strings.Contains(output, "myname") {
		t.Errorf("expected key-value output, got: %s", output)
	}
}

func TestPlainSliceOutput(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(false, true, "")
	f.Writer = &buf

	data := []testItem{{ID: 1, Name: "a"}, {ID: 2, Name: "b"}}
	if err := f.Print(data); err != nil {
		t.Fatalf("Print failed: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	// Header + 2 data rows.
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 rows), got %d: %v", len(lines), lines)
	}
}

func TestPrintMessage(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(false, false, "")
	f.Writer = &buf
	f.PrintMessage("hello world")
	if strings.TrimSpace(buf.String()) != "hello world" {
		t.Errorf("unexpected output: %s", buf.String())
	}
}

func TestPrintMessageJSON(t *testing.T) {
	var buf bytes.Buffer
	f := NewFormatter(true, false, "")
	f.Writer = &buf
	f.PrintMessage("ok")

	var result map[string]string
	if err := json.Unmarshal(buf.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["message"] != "ok" {
		t.Errorf("expected message 'ok', got %s", result["message"])
	}
}

func TestHint(t *testing.T) {
	var errBuf bytes.Buffer
	f := NewFormatter(true, false, "")
	f.ErrWriter = &errBuf
	f.Hint("hint: %s", "test")
	if !strings.Contains(errBuf.String(), "hint: test") {
		t.Errorf("expected hint on stderr, got: %s", errBuf.String())
	}
}
