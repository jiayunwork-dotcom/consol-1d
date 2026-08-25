package runbook

import (
	"errors"
	"testing"

	"consol-1d/internal/model"
)

func TestBookAddGetRemove(t *testing.T) {
	b := NewBook(4)
	e := Entry{
		ID:     "run-1",
		Input:  model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7),
		Result: model.Result{U: 0.5, SettlementRatio: 0.5},
	}
	if err := b.Add(e); err != nil {
		t.Fatal(err)
	}
	if b.Len() != 1 || b.NextSeq() != 1 {
		t.Fatalf("len=%d seq=%d", b.Len(), b.NextSeq())
	}
	got, ok := b.Get("run-1")
	if !ok || got.Input.Thickness != 10 {
		t.Fatalf("get failed: %+v", got)
	}
	if err := b.Add(e); !errors.Is(err, ErrExists) {
		t.Fatalf("duplicate err=%v", err)
	}
	if !b.Remove("run-1") {
		t.Fatal("remove failed")
	}
}

func TestBookRenameFreezeSetNote(t *testing.T) {
	b := NewBook(8)
	if err := b.Add(Entry{ID: "a", Input: model.UniformInput(1e-7, 10, model.DrainageSingle, 100, 1e7), Result: model.Result{U: 0.4, SettlementRatio: 0.4}}); err != nil {
		t.Fatal(err)
	}
	if err := b.Rename("a", "b"); err != nil {
		t.Fatal(err)
	}
	if err := b.Freeze("b"); err != nil {
		t.Fatal(err)
	}
	if err := b.SetNote("b", "changed"); !errors.Is(err, ErrFrozen) {
		t.Fatalf("frozen set note err=%v", err)
	}
}

func TestValidateRejectsBadEntries(t *testing.T) {
	bad := []Entry{
		{ID: "", Input: model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 1), Result: model.Result{U: 0.1}},
		{ID: "x", Input: model.UniformInput(0, 10, model.DrainageDouble, 100, 1), Result: model.Result{U: 0.1}},
		{ID: "x", Input: model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 1), Result: model.Result{U: 1.2}},
	}
	for i, e := range bad {
		if err := e.Validate(); err == nil {
			t.Fatalf("entry %d should fail", i)
		}
	}
}

func TestDerivedStats(t *testing.T) {
	b := NewBook(16)
	for _, e := range []Entry{
		{ID: "a", Input: model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 1e6), Result: model.Result{U: 0.2, SettlementRatio: 0.2}},
		{ID: "b", Input: model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 5e7), Result: model.Result{U: 0.5, SettlementRatio: 0.5}},
		{ID: "c", Input: model.UniformInput(1e-7, 10, model.DrainageSingle, 100, 1e9), Result: model.Result{U: 0.9, SettlementRatio: 0.9}},
	} {
		if err := b.Add(e); err != nil {
			t.Fatal(err)
		}
	}
	avg, n := b.AverageU()
	if n != 3 || avg != 0.5333333333333333 {
		t.Fatalf("avg=%v n=%d", avg, n)
	}
	maxU, id := b.MaxU()
	if id != "c" || maxU != 0.9 {
		t.Fatalf("max=%v id=%s", maxU, id)
	}
	by := b.ByDrainage()
	if by["double"] != 2 || by["single"] != 1 {
		t.Fatalf("by drainage=%v", by)
	}
	sim := b.Similar(model.UniformInput(1e-7, 10, model.DrainageDouble, 100, 1), 0.05)
	if len(sim) != 2 {
		t.Fatalf("similar=%+v", sim)
	}
	if b.NearConsolidated(0.5) != 2 {
		t.Fatalf("near=%d", b.NearConsolidated(0.5))
	}
	mean, mn := b.MeanSettlementRatio()
	if mn != 3 || mean != 0.5333333333333333 {
		t.Fatalf("mean=%v n=%d", mean, mn)
	}
}
