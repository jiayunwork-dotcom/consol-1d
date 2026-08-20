package model

import "errors"

// Sentinel errors for illegal scenario parameters. The CLI maps every
// one of them to a stderr message and a non-zero exit code, so each
// failure path below is user-visible.
var (
	// ErrCvNonPositive fires for cv <= 0.
	ErrCvNonPositive = errors.New("cv must be positive")
	// ErrThickness fires for H <= 0.
	ErrThickness = errors.New("thickness H must be positive")
	// ErrNegativeTime fires for t < 0.
	ErrNegativeTime = errors.New("time t must be non-negative")
	// ErrInitial fires for a negative initial excess pore pressure.
	ErrInitial = errors.New("initial excess pore pressure must be non-negative")
	// ErrZeroMeanPressure fires when the layer-averaged initial
	// pressure vanishes, which leaves the consolidation degree undefined.
	ErrZeroMeanPressure = errors.New("mean initial excess pore pressure must be positive")
	// ErrDrainage fires for a drainage value other than single/double.
	ErrDrainage = errors.New("drainage must be \"single\" or \"double\"")
	// ErrInitialType fires for an unknown initial_pressure.type.
	ErrInitialType = errors.New("initial_pressure.type must be \"uniform\" or \"linear\"")
	// ErrSettlement fires when mv and delta_sigma are not provided
	// together, or when a provided one is non-positive.
	ErrSettlement = errors.New("mv and delta_sigma must both be positive and provided together")
	// ErrNonFinite fires for NaN or Inf in a mandatory field.
	ErrNonFinite = errors.New("input contains NaN or Inf where a finite number is required")
	// ErrNodes fires for a profile with fewer than two nodes.
	ErrNodes = errors.New("profile needs at least 2 nodes")
)
