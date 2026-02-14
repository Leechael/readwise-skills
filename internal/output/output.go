package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"

	"github.com/itchyny/gojq"
)

// Formatter handles output formatting for CLI commands.
type Formatter struct {
	JSONMode  bool
	PlainMode bool
	JQFilter  string
	Writer    io.Writer
	ErrWriter io.Writer
}

// NewFormatter creates a new output formatter.
func NewFormatter(jsonMode, plainMode bool, jqFilter string) *Formatter {
	return &Formatter{
		JSONMode:  jsonMode,
		PlainMode: plainMode,
		JQFilter:  jqFilter,
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
	}
}

// Print outputs data in the configured format.
func (f *Formatter) Print(data interface{}) error {
	if f.JSONMode {
		return f.printJSON(data)
	}
	if f.PlainMode {
		return f.printPlain(data)
	}
	// Default to JSON if no mode specified.
	return f.printJSON(data)
}

// Hint prints a human-readable hint to stderr.
func (f *Formatter) Hint(format string, args ...interface{}) {
	fmt.Fprintf(f.ErrWriter, format+"\n", args...)
}

func (f *Formatter) printJSON(data interface{}) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}

	if f.JQFilter != "" {
		return f.applyJQ(jsonBytes)
	}

	// Pretty-print JSON.
	var pretty json.RawMessage
	if err := json.Unmarshal(jsonBytes, &pretty); err != nil {
		return err
	}
	enc := json.NewEncoder(f.Writer)
	enc.SetIndent("", "  ")
	return enc.Encode(pretty)
}

func (f *Formatter) applyJQ(jsonBytes []byte) error {
	query, err := gojq.Parse(f.JQFilter)
	if err != nil {
		return fmt.Errorf("parsing jq filter: %w", err)
	}

	var input interface{}
	if err := json.Unmarshal(jsonBytes, &input); err != nil {
		return fmt.Errorf("unmarshaling for jq: %w", err)
	}

	iter := query.Run(input)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return fmt.Errorf("jq error: %w", err)
		}
		out, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("marshaling jq result: %w", err)
		}
		// Pretty-print each result.
		var pretty json.RawMessage
		if err := json.Unmarshal(out, &pretty); err == nil {
			enc := json.NewEncoder(f.Writer)
			enc.SetIndent("", "  ")
			enc.Encode(pretty)
		} else {
			fmt.Fprintln(f.Writer, string(out))
		}
	}
	return nil
}

func (f *Formatter) printPlain(data interface{}) error {
	w := tabwriter.NewWriter(f.Writer, 0, 0, 2, ' ', 0)
	defer w.Flush()

	v := reflect.ValueOf(data)

	// Handle pointer.
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice:
		if v.Len() == 0 {
			return nil
		}
		return f.printSlicePlain(w, v)
	case reflect.Struct:
		return f.printStructPlain(w, v)
	case reflect.Map:
		return f.printMapPlain(w, v)
	default:
		fmt.Fprintln(w, fmt.Sprint(data))
	}
	return nil
}

func (f *Formatter) printSlicePlain(w *tabwriter.Writer, v reflect.Value) error {
	// Print header from first element's fields.
	first := v.Index(0)
	if first.Kind() == reflect.Ptr {
		first = first.Elem()
	}
	if first.Kind() == reflect.Struct {
		headers := structFieldNames(first)
		fmt.Fprintln(w, strings.Join(headers, "\t"))
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i)
			if elem.Kind() == reflect.Ptr {
				elem = elem.Elem()
			}
			vals := structFieldValues(elem)
			fmt.Fprintln(w, strings.Join(vals, "\t"))
		}
	} else {
		for i := 0; i < v.Len(); i++ {
			fmt.Fprintln(w, fmt.Sprint(v.Index(i).Interface()))
		}
	}
	return nil
}

func (f *Formatter) printStructPlain(w *tabwriter.Writer, v reflect.Value) error {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = field.Name
		}
		// Strip omitempty and other options.
		if idx := strings.Index(name, ","); idx != -1 {
			name = name[:idx]
		}
		fmt.Fprintf(w, "%s\t%v\n", name, formatValue(v.Field(i)))
	}
	return nil
}

func (f *Formatter) printMapPlain(w *tabwriter.Writer, v reflect.Value) error {
	for _, key := range v.MapKeys() {
		fmt.Fprintf(w, "%v\t%v\n", key.Interface(), formatValue(v.MapIndex(key)))
	}
	return nil
}

func structFieldNames(v reflect.Value) []string {
	t := v.Type()
	var names []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		name := field.Tag.Get("json")
		if name == "" || name == "-" {
			name = field.Name
		}
		if idx := strings.Index(name, ","); idx != -1 {
			name = name[:idx]
		}
		names = append(names, name)
	}
	return names
}

func structFieldValues(v reflect.Value) []string {
	t := v.Type()
	var values []string
	for i := 0; i < t.NumField(); i++ {
		if !t.Field(i).IsExported() {
			continue
		}
		values = append(values, formatValue(v.Field(i)))
	}
	return values
}

func formatValue(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		return fmt.Sprint(v.Elem().Interface())
	}
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Map {
		data, err := json.Marshal(v.Interface())
		if err != nil {
			return fmt.Sprint(v.Interface())
		}
		return string(data)
	}
	return fmt.Sprint(v.Interface())
}

// PrintMessage prints a simple string message.
func (f *Formatter) PrintMessage(msg string) {
	if f.JSONMode {
		json.NewEncoder(f.Writer).Encode(map[string]string{"message": msg})
	} else {
		fmt.Fprintln(f.Writer, msg)
	}
}
