package api

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"

	"consol-1d/internal/model"
	"consol-1d/internal/runbook"
)

type Server struct {
	mux  *http.ServeMux
	addr string
	book *runbook.Book
}

func New(addr string) *Server {
	s := &Server{
		mux:  http.NewServeMux(),
		addr: addr,
		book: runbook.NewBook(64),
	}
	s.routes()
	return s
}

func Serve(addr string) error {
	return New(addr).ListenAndServe()
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) Book() *runbook.Book {
	return s.book
}

func (s *Server) ListenAndServe() error {
	return http.ListenAndServe(s.addr, s.mux)
}

func (s *Server) routes() {
	s.mux.HandleFunc("/api/consolidate", s.handleConsolidate)
	s.mux.HandleFunc("/api/curve", s.handleCurve)
	s.mux.HandleFunc("/api/settle", s.handleSettle)
	s.mux.HandleFunc("/api/history", s.handleHistory)
	s.mux.HandleFunc("/api/health", s.handleHealth)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "GET only")
		return
	}
	writeJSON(w, http.StatusOK, s.book.List())
}

func (s *Server) handleConsolidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	var req struct {
		Input model.Input `json:"scenario"`
		Nodes int         `json:"nodes"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if req.Nodes < 2 {
		req.Nodes = 21
	}
	res, err := model.Solve(req.Input, req.Nodes)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	entry := runbook.Entry{
		ID:     fmt.Sprintf("run-%d", s.book.NextSeq()+1),
		Input:  req.Input,
		Result: res,
		Note:   string(req.Input.Drainage),
	}
	if err := s.book.Add(entry); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, consolidateResponse{
		Input:               req.Input,
		Hdr:                 res.Hdr,
		Tv:                  res.Tv,
		U:                   res.U,
		MidpointPressure:    res.MidpointPressure,
		MeanPressure:        res.MeanPressure,
		MeanInitialPressure: res.MeanInitialPressure,
		Settlement:          nullable(res.Settlement),
		UltimateSettlement:  nullable(res.UltimateSettlement),
		SettlementRatio:     res.SettlementRatio,
		Profile:             res.Profile,
		TermsUsed:           res.TermsUsed,
		RemainderBound:      res.RemainderBound,
	})
}

func (s *Server) handleCurve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	var req struct {
		Input model.Input `json:"scenario"`
		Times []float64   `json:"times"`
		Nodes int         `json:"nodes"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Times) == 0 {
		req.Times = model.DefaultCurveTimes(req.Input)
	}
	if req.Nodes < 2 {
		req.Nodes = 21
	}
	curve, err := model.ConsolidationCurve(req.Input, req.Times, req.Nodes)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	out := make([]curvePointResponse, 0, len(curve))
	for _, p := range curve {
		out = append(out, curvePointResponse{
			Time:               p.Time,
			TimeFactor:         p.TimeFactor,
			U:                  p.U,
			MidpointPressure:   p.MidpointPressure,
			Settlement:         nullable(p.Settlement),
			UltimateSettlement: nullable(p.UltimateSettlement),
		})
	}
	writeJSON(w, http.StatusOK, out)
}

func (s *Server) handleSettle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST only")
		return
	}
	var req struct {
		Input  model.Input `json:"scenario"`
		Target float64     `json:"target"`
	}
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	secs, err := model.TimeToDegree(req.Input, req.Target, 1e-6)
	if err != nil {
		writeError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"time_s": secs})
}

func readJSON(r *http.Request, v interface{}) error {
	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		return err
	}
	if len(body) == 0 {
		return fmt.Errorf("empty request body")
	}
	return json.Unmarshal(body, v)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

type consolidateResponse struct {
	Input               model.Input   `json:"scenario"`
	Hdr                 float64       `json:"hdr"`
	Tv                  float64       `json:"tv"`
	U                   float64       `json:"u"`
	MidpointPressure    float64       `json:"midpoint_pressure"`
	MeanPressure        float64       `json:"mean_pressure"`
	MeanInitialPressure float64       `json:"mean_initial_pressure"`
	Settlement          *float64      `json:"settlement,omitempty"`
	UltimateSettlement  *float64      `json:"ultimate_settlement,omitempty"`
	SettlementRatio     float64       `json:"settlement_ratio"`
	Profile             []model.Point `json:"profile"`
	TermsUsed           int           `json:"terms_used"`
	RemainderBound      float64       `json:"remainder_bound"`
}

type curvePointResponse struct {
	Time               float64  `json:"time"`
	TimeFactor         float64  `json:"time_factor"`
	U                  float64  `json:"u"`
	MidpointPressure   float64  `json:"midpoint_pressure"`
	Settlement         *float64 `json:"settlement,omitempty"`
	UltimateSettlement *float64 `json:"ultimate_settlement,omitempty"`
}

func nullable(v float64) *float64 {
	if math.IsNaN(v) || math.IsInf(v, 0) {
		return nil
	}
	return &v
}
