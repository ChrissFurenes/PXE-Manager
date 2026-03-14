package httpboot

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"PXE-Manager/internal/admin"
	"PXE-Manager/internal/boot"
	"PXE-Manager/internal/config"
	"PXE-Manager/internal/storage"
)

type Server struct {
	cfg   *config.Config
	store *storage.Store
}

func Start(cfg *config.Config, store *storage.Store) error {
	mux := http.NewServeMux()

	tmpl, err := admin.NewTemplateManager(cfg.HTTP.TemplateDir)
	if err != nil {
		return err
	}

	uploadDir := filepath.Join(cfg.HTTP.RootDir, "images")
	adminHandler := admin.NewHandler(store, tmpl, uploadDir)
	admin.RegisterRoutes(mux, adminHandler)

	s := &Server{
		cfg:   cfg,
		store: store,
	}

	mux.HandleFunc("/boot.ipxe", s.bootHandler)

	fs := http.FileServer(http.Dir(cfg.HTTP.RootDir))
	mux.Handle("/files/", http.StripPrefix("/files/", fs))

	log.Printf("[HTTP] listening on %s", cfg.HTTP.ListenAddr)
	return http.ListenAndServe(cfg.HTTP.ListenAddr, mux)
}

func (s *Server) bootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	serverBase := "http://" + s.cfg.ServerIP + s.cfg.HTTP.ListenAddr + "/files"
	mac := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("mac")))

	if mac != "" {
		client, err := s.store.GetClientByMAC(mac)
		if err == nil {
			if client.ShowMenu {
				profiles, err := s.store.ListProfiles()
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				_, _ = w.Write([]byte(boot.GenerateMenuScript(serverBase, profiles)))
				return
			}

			profile, err := s.store.GetProfile(client.ProfileID)
			if err == nil && profile.Enabled {
				_, _ = w.Write([]byte(boot.GenerateBootScript(serverBase, *profile)))
				return
			}
		}
	}

	defaultID, err := s.store.GetDefaultProfileID()
	if err == nil && defaultID > 0 {
		profile, err := s.store.GetProfile(defaultID)
		if err == nil && profile.Enabled {
			_, _ = w.Write([]byte(boot.GenerateBootScript(serverBase, *profile)))
			return
		}
	}

	profiles, err := s.store.ListProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write([]byte(boot.GenerateMenuScript(serverBase, profiles)))
}
