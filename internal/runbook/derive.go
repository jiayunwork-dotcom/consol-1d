package runbook

import (
	"sort"

	"consol-1d/internal/model"
)

func (b *Book) AverageU() (float64, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sum := 0.0
	n := 0
	for _, e := range b.items {
		sum += e.Result.U
		n++
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), n
}

func (b *Book) MaxU() (float64, string) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	best := 0.0
	id := ""
	for _, e := range b.items {
		if e.Result.U > best {
			best = e.Result.U
			id = e.ID
		}
	}
	return best, id
}

func (b *Book) ByDrainage() map[string]int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make(map[string]int)
	for _, e := range b.items {
		out[string(e.Input.Drainage)]++
	}
	return out
}

func (b *Book) Similar(target model.Input, tol float64) []Entry {
	b.mu.RLock()
	defer b.mu.RUnlock()
	out := make([]Entry, 0)
	for _, e := range b.items {
		if e.Input.Drainage != target.Drainage {
			continue
		}
		dCv := relDiff(e.Input.Cv, target.Cv)
		dH := relDiff(e.Input.Thickness, target.Thickness)
		if dCv <= tol && dH <= tol {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Seq < out[j].Seq
	})
	return out
}

func (b *Book) NearConsolidated(threshold float64) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	n := 0
	for _, e := range b.items {
		if e.Result.U >= threshold {
			n++
		}
	}
	return n
}

func (b *Book) MeanSettlementRatio() (float64, int) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	sum := 0.0
	n := 0
	for _, e := range b.items {
		sum += e.Result.SettlementRatio
		n++
	}
	if n == 0 {
		return 0, 0
	}
	return sum / float64(n), n
}

func relDiff(a, b float64) float64 {
	if b == 0 {
		return 1
	}
	d := a - b
	if d < 0 {
		d = -d
	}
	return d / b
}
