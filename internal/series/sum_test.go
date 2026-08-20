package series_test

import (
	"math"
	"testing"

	"consol-1d/internal/series"
)

// TestSumOddReciprocalSquareClosedForm checks the adaptive sum against
// the closed form sum 1/(2n+1)^2 = pi^2/8 at zero decay. This is the
// slowest family the package accepts, so it also exercises the
// remainder-bound loop at its most demanding.
func TestSumOddReciprocalSquareClosedForm(t *testing.T) {
	term := series.Term{Coefficient: 1, Power: 2, Decay: 0}
	r, err := series.Sum(term, 1e-6, 1<<25)
	if err != nil {
		t.Fatalf("Sum: expected no error, got %v", err)
	}
	want := math.Pi * math.Pi / 8
	if math.Abs(r.Sum-want) > 1e-5 {
		t.Errorf("sum of odd reciprocal squares: expected %g, got %g", want, r.Sum)
	}
	if r.RemainderBound > 1e-6 {
		t.Errorf("remainder bound: expected at most 1e-6, got %g", r.RemainderBound)
	}
	if r.TermsUsed < 1 {
		t.Errorf("terms used: expected at least 1, got %d", r.TermsUsed)
	}
}

// TestSumRespectsRemainderTolerance verifies that the returned
// remainder bound really bounds the discarded tail: the true tail from
// TermsUsed to 2*TermsUsed must fit inside the reported bound.
func TestSumRespectsRemainderTolerance(t *testing.T) {
	term := series.Term{Coefficient: 40, Power: 1, Decay: 0.5}
	r, err := series.Sum(term, 1e-8, 1<<20)
	if err != nil {
		t.Fatalf("Sum: expected no error, got %v", err)
	}
	if r.RemainderBound > 1e-8 {
		t.Errorf("remainder bound: expected at most 1e-8, got %g", r.RemainderBound)
	}
	extended := 0.0
	for i := r.TermsUsed; i < 2*r.TermsUsed; i++ {
		extended += term.Value(i)
	}
	if math.Abs(extended) > 2*r.RemainderBound {
		t.Errorf("true tail [%d, %d): expected |%g| <= 2*%g", r.TermsUsed, 2*r.TermsUsed, extended, r.RemainderBound)
	}
}

// TestRemainderBoundMonotone checks that the tail bound never grows as
// more terms are retained, for decaying and non-decaying families.
func TestRemainderBoundMonotone(t *testing.T) {
	terms := []series.Term{
		{Coefficient: 40, Power: 1, Decay: 0.5},
		{Coefficient: 1, Power: 2, Decay: 0.1},
		{Coefficient: 7, Power: 3, Decay: 0.001},
		{Coefficient: 5, Power: 2, Decay: 0},
	}
	for _, term := range terms {
		prev := term.TailBoundUpper(1)
		for n := 1; n < 60; n++ {
			cur := term.TailBoundUpper(n)
			if cur > prev {
				t.Errorf("%+v: tail bound at %d = %g grew above %g at %d", term, n, cur, prev, n-1)
			}
			prev = cur
		}
	}
}

// TestRombergMatchesClosedForm integrates simple functions with a known
// antiderivative and requires the reported error estimate to be honest.
func TestRombergMatchesClosedForm(t *testing.T) {
	v, diff, err := series.Romberg(func(x float64) float64 { return x * x * x }, 0, 1, 1e-10)
	if err != nil {
		t.Fatalf("Romberg(x^3): expected no error, got %v", err)
	}
	if want := 0.25; math.Abs(v-want) > 1e-8 {
		t.Errorf("int_0^1 x^3 dx: expected %g, got %g", want, v)
	}
	if diff > 1e-8 {
		t.Errorf("int_0^1 x^3 dx error estimate: expected at most 1e-8, got %g", diff)
	}
	v2, _, err := series.Romberg(func(x float64) float64 { return 1 / x }, 1, 2, 1e-10)
	if err != nil {
		t.Fatalf("Romberg(1/x): expected no error, got %v", err)
	}
	if want := math.Log(2); math.Abs(v2-want) > 1e-8 {
		t.Errorf("int_1^2 1/x dx: expected %g, got %g", want, v2)
	}
}

// TestStricterToleranceAddsTerms checks that demanding a smaller tail
// makes the adaptive loop keep more terms.
func TestStricterToleranceAddsTerms(t *testing.T) {
	term := series.Term{Coefficient: 1, Power: 1, Decay: 0.05}
	lax, err := series.Sum(term, 1e-6, 1<<25)
	if err != nil {
		t.Fatalf("Sum(1e-6): expected no error, got %v", err)
	}
	strict, err := series.Sum(term, 1e-12, 1<<25)
	if err != nil {
		t.Fatalf("Sum(1e-12): expected no error, got %v", err)
	}
	if strict.TermsUsed <= lax.TermsUsed {
		t.Errorf("strict tolerance should keep more terms: lax %d, strict %d", lax.TermsUsed, strict.TermsUsed)
	}
	if strict.RemainderBound > 1e-12 {
		t.Errorf("strict remainder bound: expected at most 1e-12, got %g", strict.RemainderBound)
	}
}

// TestSumRejectsNonsummableFamily verifies that a harmonic-like family
// (Power 1, zero decay) is refused instead of silently diverging.
func TestSumRejectsNonsummableFamily(t *testing.T) {
	term := series.Term{Coefficient: 3, Power: 1, Decay: 0}
	if _, err := series.Sum(term, 1e-9, 100); err == nil {
		t.Error("Sum of a power-1 zero-decay family: expected an error, got nil")
	}
}
