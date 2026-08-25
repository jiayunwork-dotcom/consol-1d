package runbook

import (
	"os"
	"path/filepath"
	"testing"

	"consol-1d/internal/model"
)

func TestSaveLoadRoundTrip(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "book.json")
	b := NewBook(12)
	if err := b.Add(Entry{
		ID:     "a",
		Input:  model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7),
		Result: model.Result{U: 0.5, SettlementRatio: 0.5},
	}); err != nil {
		t.Fatal(err)
	}
	if err := b.SaveFile(path); err != nil {
		t.Fatal(err)
	}
	loaded, err := LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Len() != 1 || loaded.NextSeq() != 1 {
		t.Fatalf("loaded len=%d seq=%d", loaded.Len(), loaded.NextSeq())
	}
	got, ok := loaded.Get("a")
	if !ok || got.Result.U != 0.5 {
		t.Fatalf("loaded mismatch: %+v", got)
	}
}

func TestLoadRejectsCorruptSnapshot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte(`{"version":7,"max":4,"seq":1,"entries":[]}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadFile(path); err == nil {
		t.Fatal("bad version should fail")
	}
}

func TestImportRebuildsBook(t *testing.T) {
	b := NewBook(10)
	if err := b.Add(Entry{
		ID:     "x",
		Input:  model.UniformInput(1e-7, 10, model.DrainageSingle, 100, 1e7),
		Result: model.Result{U: 0.4, SettlementRatio: 0.4},
	}); err != nil {
		t.Fatal(err)
	}
	other := NewBook(3)
	if err := other.Import(b.Export()); err != nil {
		t.Fatal(err)
	}
	if other.Len() != 1 {
		t.Fatalf("import len=%d", other.Len())
	}
}
