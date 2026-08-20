package model

import (
	"fmt"
	"math"
)

// finite reports whether v is a usable real number.
func finite(v float64) bool {
	return !math.IsNaN(v) && !math.IsInf(v, 0)
}

// Validate rejects every illegal scenario before any series is built:
// cv <= 0, H <= 0, t < 0, a negative initial pressure, an unknown
// drainage path or initial-profile type, a vanishing mean initial
// pressure, and a half-provided settlement pair. NaN or Inf anywhere in
// a mandatory field is also rejected, because a NaN silently poisons
// every downstream series.
func (in Input) Validate() error {
	switch in.Drainage {
	case DrainageSingle, DrainageDouble:
	default:
		return fmt.Errorf("%w (got %q)", ErrDrainage, in.Drainage)
	}
	if !finite(in.Cv) || !finite(in.Thickness) || !finite(in.Time) {
		return ErrNonFinite
	}
	if in.Cv <= 0 {
		return commitV(fmt.Errorf("%w (got cv=%g)", ErrCvNonPositive, in.Cv))
	}
	if in.Thickness <= 0 {
		return fmt.Errorf("%w (got H=%g)", ErrThickness, in.Thickness)
	}
	if in.Time < 0 {
		return fmt.Errorf("%w (got t=%g)", ErrNegativeTime, in.Time)
	}
	if err := in.Initial.validate(); err != nil {
		return err
	}
	if (in.Mv == nil) != (in.DeltaSigma == nil) {
		return fmt.Errorf("%w (mv=%v, delta_sigma=%v)", ErrSettlement, in.Mv, in.DeltaSigma)
	}
	if in.Mv != nil {
		if !finite(*in.Mv) || *in.Mv <= 0 {
			return fmt.Errorf("%w (got mv=%g)", ErrSettlement, *in.Mv)
		}
		if !finite(*in.DeltaSigma) || *in.DeltaSigma <= 0 {
			return fmt.Errorf("%w (got delta_sigma=%g)", ErrSettlement, *in.DeltaSigma)
		}
	}
	return nil
}

// validate checks the initial pressure profile alone.
func (ip InitialPressure) validate() error {
	switch ip.Type {
	case InitialUniform:
		if !finite(ip.U0) {
			return fmt.Errorf("%w (u0=%g)", ErrNonFinite, ip.U0)
		}
		if ip.U0 < 0 {
			return fmt.Errorf("%w (u0=%g)", ErrInitial, ip.U0)
		}
		if ip.U0 == 0 {
			return ErrZeroMeanPressure
		}
	case InitialLinear:
		if !finite(ip.UA) || !finite(ip.UB) {
			return fmt.Errorf("%w (ua=%g, ub=%g)", ErrNonFinite, ip.UA, ip.UB)
		}
		if ip.UA < 0 || ip.UB < 0 {
			return fmt.Errorf("%w (ua=%g, ub=%g)", ErrInitial, ip.UA, ip.UB)
		}
		if ip.UA+ip.UB == 0 {
			return ErrZeroMeanPressure
		}
	default:
		return fmt.Errorf("%w (got %q)", ErrInitialType, ip.Type)
	}
	return nil
}
