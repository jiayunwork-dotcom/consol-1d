package model

type Drainage string

const (
	DrainageSingle Drainage = "single"
	DrainageDouble Drainage = "double"
)

type InitialKind string

const (
	InitialUniform InitialKind = "uniform"
	InitialLinear  InitialKind = "linear"
)

type InitialPressure struct {
	Type InitialKind `json:"type"`
	U0   float64     `json:"u0"`
	UA   float64     `json:"ua"`
	UB   float64     `json:"ub"`
}

type Input struct {
	Cv         float64         `json:"cv"`
	Thickness  float64         `json:"thickness"`
	Drainage   Drainage        `json:"drainage"`
	Initial    InitialPressure `json:"initial_pressure"`
	Time       float64         `json:"time"`
	Mv         *float64        `json:"mv"`
	DeltaSigma *float64        `json:"delta_sigma"`
}
