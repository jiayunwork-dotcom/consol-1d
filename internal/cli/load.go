// Package cli wires the JSON scenario files to the consolidation solver
// and the report package. It owns the command-line surface: flag
// parsing, strict JSON decoding and the non-zero exit behaviour that
// makes every illegal input visible.
package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"consol-1d/internal/model"
)

// LoadInput reads a scenario JSON file. Unknown fields are rejected so
// a typo like "thicknes" cannot silently keep the default thickness and
// corrupt the physics; the error names the file so the user can fix the
// right one.
func LoadInput(path string) (model.Input, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Input{}, fmt.Errorf("read scenario %s: %w", path, err)
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var in model.Input
	if err := dec.Decode(&in); err != nil {
		return model.Input{}, fmt.Errorf("parse scenario %s: %w", path, err)
	}
	return in, nil
}

// ParseInput is the same decoder as LoadInput but operates on a byte
// slice; it exists so tests and embedders can feed an in-memory
// scenario without touching the filesystem.
func ParseInput(data []byte) (model.Input, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var in model.Input
	if err := dec.Decode(&in); err != nil {
		return model.Input{}, err
	}
	return in, nil
}
