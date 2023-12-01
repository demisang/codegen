package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/demisang/codegen/internal"
)

func (s *Server) templates(w http.ResponseWriter, r *http.Request) {
	templates, err := s.app.GetTemplatesList(r.Context())
	if err != nil {
		errResponse(w, r, 500, err)
		return
	}

	response(w, r, http.StatusOK, templates)
}

func (s *Server) rawList(w http.ResponseWriter, r *http.Request) {
	var requestParams internal.ReplaceOptions

	err := json.NewDecoder(r.Body).Decode(&requestParams)
	if err != nil {
		errResponse(w, r, 400, fmt.Errorf("request params format: %w", err))
		return
	}

	templates, err := s.app.RawList(r.Context(), requestParams)
	if err != nil {
		errResponse(w, r, 500, err)
		return
	}

	response(w, r, http.StatusOK, templates)
}

func (s *Server) previewList(w http.ResponseWriter, r *http.Request) {
	var requestParams internal.ReplaceOptions

	err := json.NewDecoder(r.Body).Decode(&requestParams)
	if err != nil {
		errResponse(w, r, 400, fmt.Errorf("request params format: %w", err))
		return
	}

	templates, err := s.app.PreviewList(r.Context(), requestParams)
	if err != nil {
		errResponse(w, r, 500, err)
		return
	}

	response(w, r, http.StatusOK, templates)
}

func (s *Server) generate(w http.ResponseWriter, r *http.Request) {
	var requestParams internal.ReplaceOptions

	err := json.NewDecoder(r.Body).Decode(&requestParams)
	if err != nil {
		errResponse(w, r, 400, fmt.Errorf("request params format: %w", err))
		return
	}

	helpMessage, err := s.app.Generate(r.Context(), requestParams)
	if err != nil {
		errResponse(w, r, 500, err)
		return
	}

	response(w, r, http.StatusOK, helpMessage)
}

func (s *Server) directories(w http.ResponseWriter, r *http.Request) {
	selectedDir := r.URL.Query().Get("selected")

	dirs, err := s.app.GetDirectories(r.Context(), selectedDir)
	if err != nil {
		errResponse(w, r, 500, err)
		return
	}

	response(w, r, http.StatusOK, dirs)
}
