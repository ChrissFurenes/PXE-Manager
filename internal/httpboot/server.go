package httpboot

import (
	"log"
	"path/filepath"
	"strings"

	"PXE-Manager/internal/admin"
	"PXE-Manager/internal/boot"
	"PXE-Manager/internal/config"
	"PXE-Manager/internal/storage"

	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg   *config.Config
	store *storage.Store
}

func Start(cfg *config.Config, store *storage.Store) error {
	router := gin.Default()

	router.LoadHTMLGlob(filepath.Join(cfg.HTTP.TemplateDir, "*.html"))

	uploadDir := filepath.Join(cfg.HTTP.RootDir, "images")

	adminHandler := admin.NewHandler(store, uploadDir)
	admin.RegisterRoutes(router, adminHandler)

	s := &Server{
		cfg:   cfg,
		store: store,
	}

	router.GET("/boot.ipxe", s.bootHandler)
	router.Static("/files", cfg.HTTP.RootDir)

	log.Printf("[HTTP] listening on %s", cfg.HTTP.ListenAddr)
	return router.Run(cfg.HTTP.ListenAddr)
}

func (s *Server) bootHandler(c *gin.Context) {
	serverBase := "http://" + s.cfg.ServerIP + s.cfg.HTTP.ListenAddr + "/files"
	mac := strings.TrimSpace(strings.ToLower(c.Query("mac")))

	if mac != "" {
		client, err := s.store.GetClientByMAC(mac)
		if err == nil {
			if client.ShowMenu {
				profiles, err := s.store.ListProfiles()
				if err != nil {
					c.String(500, err.Error())
					return
				}
				c.Data(200, "text/plain; charset=utf-8", []byte(boot.GenerateMenuScript(serverBase, profiles)))
				return
			}

			profile, err := s.store.GetProfile(client.ProfileID)
			if err == nil && profile.Enabled {
				c.Data(200, "text/plain; charset=utf-8", []byte(boot.GenerateBootScript(serverBase, *profile)))
				return
			}
		}
	}

	defaultID, err := s.store.GetDefaultProfileID()
	if err == nil && defaultID > 0 {
		profile, err := s.store.GetProfile(defaultID)
		if err == nil && profile.Enabled {
			c.Data(200, "text/plain; charset=utf-8", []byte(boot.GenerateBootScript(serverBase, *profile)))
			return
		}
	}

	profiles, err := s.store.ListProfiles()
	if err != nil {
		c.String(500, err.Error())
		return
	}

	c.Data(200, "text/plain; charset=utf-8", []byte(boot.GenerateMenuScript(serverBase, profiles)))
}
