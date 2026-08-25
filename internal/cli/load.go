package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"consol-1d/internal/model"
)

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

func ParseInput(data []byte) (model.Input, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var in model.Input
	if err := dec.Decode(&in); err != nil {
		return model.Input{}, err
	}
	return in, nil
}
