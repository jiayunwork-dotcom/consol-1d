package model

import "errors"

var (
	ErrCvNonPositive    = errors.New("cv must be positive")
	ErrThickness        = errors.New("thickness H must be positive")
	ErrNegativeTime     = errors.New("time t must be non-negative")
	ErrInitial          = errors.New("initial excess pore pressure must be non-negative")
	ErrZeroMeanPressure = errors.New("mean initial excess pore pressure must be positive")
	ErrDrainage         = errors.New("drainage must be \"single\" or \"double\"")
	ErrInitialType      = errors.New("initial_pressure.type must be \"uniform\" or \"linear\"")
	ErrSettlement       = errors.New("mv and delta_sigma must both be positive and provided together")
	ErrNonFinite        = errors.New("input contains NaN or Inf where a finite number is required")
	ErrNodes            = errors.New("profile needs at least 2 nodes")
)
