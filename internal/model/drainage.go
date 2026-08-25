package model

import "math"

type DrainagePath struct {
	Kind Drainage
	H    float64
	Hdr  float64
}

func NewDrainagePath(d Drainage, h float64) DrainagePath {
	hdr := h
	if d == DrainageDouble {
		hdr = h / 2
	}
	return DrainagePath{Kind: d, H: h, Hdr: hdr}
}

func (p DrainagePath) TimeFactor(cv, t float64) float64 {
	return cv * t / (p.Hdr * p.Hdr)
}

func (p DrainagePath) Drained(fraction float64) bool {
	switch p.Kind {
	case DrainageSingle:
		return fraction == 0
	default:
		return fraction == 0 || fraction == 1
	}
}

func (p DrainagePath) Eigenvalue(n int) float64 {
	k := float64(2*n + 1)
	if p.Kind == DrainageSingle {
		return k * math.Pi / (2 * p.H)
	}
	return k * math.Pi / p.H
}

func (p DrainagePath) EigenIntegral(n int) float64 {
	lam := p.Eigenvalue(n)
	return (1 - math.Cos(lam*p.H)) / lam
}
