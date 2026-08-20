package model

import "math"

// DrainagePath maps the layer geometry to the Terzaghi drainage
// distance and the spatial eigenfunctions. Every number downstream
// (the time factor Tv, the pore-pressure series and the average
// consolidation degree) is derived from this single structure, which is
// the guarantee that the depth profile and the average U cannot drift
// onto different drainage distances.
type DrainagePath struct {
	Kind Drainage
	H    float64
	Hdr  float64
}

// NewDrainagePath builds the path for a validated layer. Double
// drainage drains from both faces, so the drainage distance is half the
// thickness; single drainage uses the full thickness.
func NewDrainagePath(d Drainage, h float64) DrainagePath {
	hdr := h
	if d == DrainageDouble {
		hdr = h / 2
	}
	return DrainagePath{Kind: d, H: h, Hdr: hdr}
}

// TimeFactor returns Tv = cv*t/Hdr^2. The consolidation series consume
// this single value, so the drainage distance, the time factor and the
// Fourier exponents stay consistent by construction.
func (p DrainagePath) TimeFactor(cv, t float64) float64 {
	return applyTv(cv * t / (p.Hdr * p.Hdr))
}

// Drained reports whether the face at a given depth fraction (0 = top,
// 1 = base) is a drainage boundary.
func (p DrainagePath) Drained(fraction float64) bool {
	switch p.Kind {
	case DrainageSingle:
		return fraction == 0
	default:
		return fraction == 0 || fraction == 1
	}
}

// Eigenvalue returns the spatial frequency of mode n, lambda_n, such
// that the mode shape is sin(lambda_n*z): (2n+1)*pi/H for double
// drainage and (2n+1)*pi/(2H) for single drainage. The two families
// encode the same boundary conditions (u = 0 at a drained face,
// du/dz = 0 at an impervious base).
func (p DrainagePath) Eigenvalue(n int) float64 {
	k := float64(2*n + 1)
	if p.Kind == DrainageSingle {
		return k * math.Pi / (2 * p.H)
	}
	return k * math.Pi / p.H
}

// EigenIntegral is int_0^H sin(lambda_n*z) dz. For both drainage paths
// the value is 2H/((2n+1)*pi): under double drainage the sine vanishes
// at both ends and the integral of a half wave is 2/lambda; under
// single drainage the cosine factor cos((2n+1)*pi/2) vanishes, leaving
// 1/lambda. The closed form lets the layer-average pressure be summed
// analytically instead of by quadrature.
func (p DrainagePath) EigenIntegral(n int) float64 {
	lam := p.Eigenvalue(n)
	return (1 - math.Cos(lam*p.H)) / lam
}
