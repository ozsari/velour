package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ozsari/velour/internal/models"
	"github.com/ozsari/velour/internal/services"
)

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, status int, message string) {
	jsonResponse(w, status, map[string]string{"error": message})
}

// isNative returns true if the system is configured for native installs
func (s *Server) isNative() bool {
	return s.cfg.InstallMode == "native"
}

func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")
		if header == "" {
			jsonError(w, http.StatusUnauthorized, "missing authorization header")
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		if token == header {
			jsonError(w, http.StatusUnauthorized, "invalid authorization format")
			return
		}

		user, err := s.auth.ValidateToken(token)
		if err != nil {
			jsonError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		r.Header.Set("X-User-ID", fmt.Sprintf("%d", user.ID))
		r.Header.Set("X-Username", user.Username)
		next(w, r)
	}
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"version": s.cfg.Version,
	})
}

func (s *Server) handleConfig(w http.ResponseWriter, r *http.Request) {
	jsonResponse(w, http.StatusOK, map[string]string{
		"version":      s.cfg.Version,
		"install_mode": s.cfg.InstallMode,
	})
}

func (s *Server) handleSetupStatus(w http.ResponseWriter, r *http.Request) {
	needsSetup, err := s.auth.NeedsSetup()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to check setup status")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]bool{"needs_setup": needsSetup})
}

func (s *Server) handleSetup(w http.ResponseWriter, r *http.Request) {
	needsSetup, _ := s.auth.NeedsSetup()
	if !needsSetup {
		jsonError(w, http.StatusBadRequest, "setup already completed")
		return
	}

	var req models.SetupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Username == "" || req.Password == "" {
		jsonError(w, http.StatusBadRequest, "username and password required")
		return
	}

	if len(req.Password) < 6 {
		jsonError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	user, err := s.auth.CreateUser(req.Username, req.Password, true)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	// Save credentials for app provisioning (like swizzin)
	s.cfg.AppUsername = req.Username
	s.cfg.AppPassword = req.Password
	if err := s.cfg.Save(); err != nil {
		log.Printf("Warning: failed to save app credentials: %v", err)
	}
	if s.native != nil {
		s.native.SetCredentials(req.Username, req.Password)
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "setup completed",
		"user":    user,
	})
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, user, err := s.auth.Login(req.Username, req.Password)
	if err != nil {
		jsonError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	jsonResponse(w, http.StatusOK, models.LoginResponse{
		Token: token,
		User:  *user,
	})
}

func (s *Server) handleSystemInfo(w http.ResponseWriter, r *http.Request) {
	info, err := s.monitor.GetSystemInfo()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to get system info")
		return
	}
	jsonResponse(w, http.StatusOK, info)
}

func (s *Server) handleListServices(w http.ResponseWriter, r *http.Request) {
	var allServices []models.Service

	if s.isNative() {
		// Native mode: list systemd services
		if s.native != nil {
			native, err := s.native.ListManaged(r.Context())
			if err != nil {
				log.Printf("Failed to list native services: %v", err)
			} else {
				allServices = append(allServices, native...)
			}
		}
	} else {
		// Docker mode: list containers
		if s.docker != nil {
			managed, err := s.docker.ListManaged(r.Context())
			if err != nil {
				log.Printf("Failed to list Docker services: %v", err)
			} else {
				allServices = append(allServices, managed...)
			}
		}
	}

	if allServices == nil {
		allServices = []models.Service{}
	}

	jsonResponse(w, http.StatusOK, allServices)
}

func (s *Server) handleServiceCatalog(w http.ResponseWriter, r *http.Request) {
	catalog := services.GetRegistry()
	jsonResponse(w, http.StatusOK, catalog)
}

func (s *Server) handleInstallService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	def := services.FindByID(id)
	if def == nil {
		jsonError(w, http.StatusNotFound, "service not found in catalog")
		return
	}

	if s.isNative() {
		if s.native == nil {
			jsonError(w, http.StatusServiceUnavailable, "native installer not available")
			return
		}
		if def.Native == nil {
			jsonError(w, http.StatusBadRequest, "native install not supported for this service")
			return
		}
		go func() {
			ctx := context.Background()
			if err := s.native.Install(ctx, def); err != nil {
				log.Printf("Failed to install %s (native): %v", id, err)
			} else {
				log.Printf("Successfully installed %s (native)", id)
			}
		}()
	} else {
		if s.docker == nil {
			jsonError(w, http.StatusServiceUnavailable, "Docker not available")
			return
		}
		go func() {
			ctx := context.Background()
			if err := s.docker.Install(ctx, def); err != nil {
				log.Printf("Failed to install %s (docker): %v", id, err)
			} else {
				log.Printf("Successfully installed %s (docker)", id)
			}
		}()
	}

	jsonResponse(w, http.StatusAccepted, map[string]string{
		"message": "installation started",
		"service": id,
		"type":    s.cfg.InstallMode,
	})
}

func (s *Server) handleStartService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var err error
	if s.isNative() {
		err = s.native.Start(r.Context(), id)
	} else {
		if s.docker == nil {
			jsonError(w, http.StatusServiceUnavailable, "Docker not available")
			return
		}
		err = s.docker.Start(r.Context(), id)
	}

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to start service")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "started"})
}

func (s *Server) handleStopService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var err error
	if s.isNative() {
		err = s.native.Stop(r.Context(), id)
	} else {
		if s.docker == nil {
			jsonError(w, http.StatusServiceUnavailable, "Docker not available")
			return
		}
		err = s.docker.Stop(r.Context(), id)
	}

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to stop service")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "stopped"})
}

func (s *Server) handleRestartService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var err error
	if s.isNative() {
		err = s.native.Restart(r.Context(), id)
	} else {
		if s.docker == nil {
			jsonError(w, http.StatusServiceUnavailable, "Docker not available")
			return
		}
		err = s.docker.Restart(r.Context(), id)
	}

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to restart service")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "restarted"})
}

// ── Network Stats ──

func (s *Server) handleNetworkMonths(w http.ResponseWriter, r *http.Request) {
	months, err := s.netTracker.ListMonths()
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to list network months")
		return
	}
	if months == nil {
		months = []models.MonthlyNetStats{}
	}
	jsonResponse(w, http.StatusOK, months)
}

func (s *Server) handleNetworkMonth(w http.ResponseWriter, r *http.Request) {
	month := r.PathValue("month")
	stats, err := s.netTracker.GetMonth(month)
	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to get network stats")
		return
	}
	jsonResponse(w, http.StatusOK, stats)
}

func (s *Server) handleRemoveService(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var err error
	if s.isNative() {
		err = s.native.Remove(r.Context(), id)
	} else {
		if s.docker == nil {
			jsonError(w, http.StatusServiceUnavailable, "Docker not available")
			return
		}
		err = s.docker.Remove(r.Context(), id)
	}

	if err != nil {
		jsonError(w, http.StatusInternalServerError, "failed to remove service")
		return
	}
	jsonResponse(w, http.StatusOK, map[string]string{"message": "removed"})
}
