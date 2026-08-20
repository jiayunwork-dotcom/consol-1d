package model

// Drainage names the drainage path of the clay layer.
type Drainage string

const (
	// DrainageSingle drains from one face at z = 0; the base z = H is
	// impervious, so the drainage distance is the full thickness.
	DrainageSingle Drainage = "single"
	// DrainageDouble drains from both faces z = 0 and z = H, so the
	// drainage distance is half the thickness.
	DrainageDouble Drainage = "double"
)

// InitialKind names the shape of the initial excess pore-pressure
// profile u(z, 0).
type InitialKind string

const (
	// InitialUniform is a constant profile u(z, 0) = u0.
	InitialUniform InitialKind = "uniform"
	// InitialLinear is a straight line from ua at the top face z = 0
	// to ub at z = H.
	InitialLinear InitialKind = "linear"
)

// InitialPressure is the excess pore pressure field at t = 0. A
// uniform profile needs only U0; a linear profile interpolates between
// UA and UB. The uniform profile is the ua = ub = u0 special case of
// the linear one, and the solver preserves that symmetry.
type InitialPressure struct {
	Type InitialKind `json:"type"`
	U0   float64     `json:"u0"`
	UA   float64     `json:"ua"`
	UB   float64     `json:"ub"`
}

// Input is the full JSON scenario for one consolidation check. All
// pressures are in kPa, lengths in metres and time in seconds, so cv is
// in m^2/s and mv in 1/kPa (equivalently m^2/kN).
type Input struct {
	Cv         float64         `json:"cv"`
	Thickness  float64         `json:"thickness"`
	Drainage   Drainage        `json:"drainage"`
	Initial    InitialPressure `json:"initial_pressure"`
	Time       float64         `json:"time"`
	Mv         *float64        `json:"mv"`
	DeltaSigma *float64        `json:"delta_sigma"`
}
