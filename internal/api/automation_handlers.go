package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ozsari/velour/internal/models"
)

func (s *Server) handleListRules(w http.ResponseWriter, r *http.Request) {
	if s.automation == nil {
		jsonError(w, http.StatusServiceUnavailable, "automation not available")
		return
	}
	rules, err := s.automation.ListRules()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list rules")
		return
	}
	if rules == nil {
		rules = []models.AutomationRule{}
	}
	jsonResponse(w, http.StatusOK, rules)
}

func (s *Server) handleCreateRule(w http.ResponseWriter, r *http.Request) {
	if s.automation == nil {
		jsonError(w, http.StatusServiceUnavailable, "automation not available")
		return
	}

	var req models.CreateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		jsonError(w, http.StatusBadRequest, "name is required")
		return
	}

	now := time.Now()
	rule := models.AutomationRule{
		ID:        fmt.Sprintf("rule_%d", now.UnixNano()),
		Name:      req.Name,
		Enabled:   true,
		Trigger:   req.Trigger,
		Action:    req.Action,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.automation.CreateRule(rule); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to create rule")
		return
	}

	jsonResponse(w, http.StatusCreated, rule)
}

func (s *Server) handleGetRule(w http.ResponseWriter, r *http.Request) {
	if s.automation == nil {
		jsonError(w, http.StatusServiceUnavailable, "automation not available")
		return
	}

	id := r.PathValue("id")
	rule, err := s.automation.GetRule(id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "rule not found")
		return
	}
	jsonResponse(w, http.StatusOK, rule)
}

func (s *Server) handleUpdateRule(w http.ResponseWriter, r *http.Request) {
	if s.automation == nil {
		jsonError(w, http.StatusServiceUnavailable, "automation not available")
		return
	}

	id := r.PathValue("id")
	existing, err := s.automation.GetRule(id)
	if err != nil {
		jsonError(w, http.StatusNotFound, "rule not found")
		return
	}

	var req models.CreateRuleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	existing.Name = req.Name
	existing.Trigger = req.Trigger
	existing.Action = req.Action
	existing.UpdatedAt = time.Now()

	if err := s.automation.UpdateRule(*existing); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to update rule")
		return
	}

	jsonResponse(w, http.StatusOK, existing)
}

func (s *Server) handleDeleteRule(w http.ResponseWriter, r *http.Request) {
	if s.automation == nil {
		jsonError(w, http.StatusServiceUnavailable, "automation not available")
		return
	}

	id := r.PathValue("id")
	if err := s.automation.DeleteRule(id); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to delete rule")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "deleted"})
}

func (s *Server) handleToggleRule(w http.ResponseWriter, r *http.Request) {
	if s.automation == nil {
		jsonError(w, http.StatusServiceUnavailable, "automation not available")
		return
	}

	id := r.PathValue("id")
	var req struct {
		Enabled bool `json:"enabled"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := s.automation.ToggleRule(id, req.Enabled); err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to toggle rule")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "toggled"})
}

func (s *Server) handleListTemplates(w http.ResponseWriter, r *http.Request) {
	if s.automation == nil {
		jsonError(w, http.StatusServiceUnavailable, "automation not available")
		return
	}
	jsonResponse(w, http.StatusOK, s.automation.GetTemplates())
}
