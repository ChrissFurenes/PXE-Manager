package admin

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"PXE-Manager/internal/storage"
)

type Handler struct {
	store     *storage.Store
	tmpl      *TemplateManager
	uploadDir string
}

func NewHandler(store *storage.Store, tmpl *TemplateManager, uploadDir string) *Handler {
	return &Handler{
		store:     store,
		tmpl:      tmpl,
		uploadDir: uploadDir,
	}
}

func (h *Handler) Dashboard(w http.ResponseWriter, r *http.Request) {
	profiles, _ := h.store.ListProfiles()
	clients, _ := h.store.ListClients()
	assets, _ := h.store.ListAssets()
	defaultProfileID, _ := h.store.GetDefaultProfileID()

	data := map[string]any{
		"ProfilesCount":    len(profiles),
		"ClientsCount":     len(clients),
		"AssetsCount":      len(assets),
		"DefaultProfileID": defaultProfileID,
	}

	h.tmpl.Render(w, "dashboard.html", data)
}

func (h *Handler) Profiles(w http.ResponseWriter, r *http.Request) {
	profiles, err := h.store.ListProfiles()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defaultProfileID, _ := h.store.GetDefaultProfileID()

	data := map[string]any{
		"Profiles":         profiles,
		"DefaultProfileID": defaultProfileID,
	}

	h.tmpl.Render(w, "profiles.html", data)
}

func (h *Handler) NewProfile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		assets, _ := h.store.ListAssets()
		h.tmpl.Render(w, "profile_form.html", map[string]any{
			"Assets":  assets,
			"Profile": nil,
			"IsEdit":  false,
		})
		return

	case http.MethodPost:
		enabled := r.FormValue("enabled") == "on"

		p := &storage.Profile{
			Name:      strings.TrimSpace(r.FormValue("name")),
			BootMode:  strings.TrimSpace(r.FormValue("boot_mode")),
			BootType:  strings.TrimSpace(r.FormValue("boot_type")),
			Kernel:    strings.TrimSpace(r.FormValue("kernel")),
			Initrd:    strings.TrimSpace(r.FormValue("initrd")),
			ImagePath: strings.TrimSpace(r.FormValue("image_path")),
			Cmdline:   strings.TrimSpace(r.FormValue("cmdline")),
			Enabled:   enabled,
		}

		if err := h.store.CreateProfile(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/profiles", http.StatusSeeOther)
		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) EditProfile(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid profile id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		profile, err := h.store.GetProfile(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		assets, _ := h.store.ListAssets()

		h.tmpl.Render(w, "profile_form.html", map[string]any{
			"Assets":  assets,
			"Profile": profile,
			"IsEdit":  true,
		})
		return

	case http.MethodPost:
		enabled := r.FormValue("enabled") == "on"

		p := &storage.Profile{
			ID:        id,
			Name:      strings.TrimSpace(r.FormValue("name")),
			BootMode:  strings.TrimSpace(r.FormValue("boot_mode")),
			BootType:  strings.TrimSpace(r.FormValue("boot_type")),
			Kernel:    strings.TrimSpace(r.FormValue("kernel")),
			Initrd:    strings.TrimSpace(r.FormValue("initrd")),
			ImagePath: strings.TrimSpace(r.FormValue("image_path")),
			Cmdline:   strings.TrimSpace(r.FormValue("cmdline")),
			Enabled:   enabled,
		}

		if err := h.store.UpdateProfile(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/profiles", http.StatusSeeOther)
		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) DeleteProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid profile id", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteProfile(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/profiles", http.StatusSeeOther)
}

func (h *Handler) Clients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.store.ListClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.tmpl.Render(w, "clients.html", map[string]any{
		"Clients": clients,
	})
}

func (h *Handler) NewClient(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		profiles, _ := h.store.ListProfiles()
		h.tmpl.Render(w, "client_form.html", map[string]any{
			"Profiles": profiles,
			"Client":   nil,
			"IsEdit":   false,
		})
		return

	case http.MethodPost:
		showMenu := r.FormValue("show_menu") == "on"
		profileID, _ := strconv.ParseInt(r.FormValue("profile_id"), 10, 64)

		c := &storage.Client{
			MAC:         strings.TrimSpace(strings.ToLower(r.FormValue("mac"))),
			Hostname:    strings.TrimSpace(r.FormValue("hostname")),
			ProfileID:   profileID,
			ShowMenu:    showMenu,
			Description: strings.TrimSpace(r.FormValue("description")),
		}

		if err := h.store.CreateClient(c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/clients", http.StatusSeeOther)
		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) EditClient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid client id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		client, err := h.store.GetClient(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		profiles, _ := h.store.ListProfiles()

		h.tmpl.Render(w, "client_form.html", map[string]any{
			"Profiles": profiles,
			"Client":   client,
			"IsEdit":   true,
		})
		return

	case http.MethodPost:
		showMenu := r.FormValue("show_menu") == "on"
		profileID, _ := strconv.ParseInt(r.FormValue("profile_id"), 10, 64)

		c := &storage.Client{
			ID:          id,
			MAC:         strings.TrimSpace(strings.ToLower(r.FormValue("mac"))),
			Hostname:    strings.TrimSpace(r.FormValue("hostname")),
			ProfileID:   profileID,
			ShowMenu:    showMenu,
			Description: strings.TrimSpace(r.FormValue("description")),
		}

		if err := h.store.UpdateClient(c); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/clients", http.StatusSeeOther)
		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) DeleteClient(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid client id", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteClient(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/clients", http.StatusSeeOther)
}

func (h *Handler) SetDefaultProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.FormValue("profile_id"), 10, 64)
	if err != nil {
		http.Error(w, "invalid profile id", http.StatusBadRequest)
		return
	}

	if err := h.store.SetDefaultProfileID(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/profiles", http.StatusSeeOther)
}

func (h *Handler) Assets(w http.ResponseWriter, r *http.Request) {
	assets, err := h.store.ListAssets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.tmpl.Render(w, "assets.html", map[string]any{
		"Assets": assets,
	})
}

func (h *Handler) NewAsset(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.tmpl.Render(w, "asset_form.html", map[string]any{
			"Asset":  nil,
			"IsEdit": false,
		})
		return

	case http.MethodPost:
		err := r.ParseMultipartForm(10 << 30)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		name := strings.TrimSpace(r.FormValue("name"))
		description := strings.TrimSpace(r.FormValue("description"))
		fileType := strings.TrimSpace(r.FormValue("file_type"))

		if name == "" {
			name = header.Filename
		}

		if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		safeName := filepath.Base(header.Filename)
		dstPath := filepath.Join(h.uploadDir, safeName)

		dst, err := os.Create(dstPath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		written, err := io.Copy(dst, file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		asset := &storage.Asset{
			Name:        name,
			FileName:    safeName,
			FilePath:    "/images/" + safeName,
			FileType:    fileType,
			SizeBytes:   written,
			Description: description,
		}

		if err := h.store.CreateAsset(asset); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/assets", http.StatusSeeOther)
		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) EditAsset(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid asset id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		asset, err := h.store.GetAsset(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		h.tmpl.Render(w, "asset_form.html", map[string]any{
			"Asset":  asset,
			"IsEdit": true,
		})
		return

	case http.MethodPost:
		asset := &storage.Asset{
			ID:          id,
			Name:        strings.TrimSpace(r.FormValue("name")),
			FileType:    strings.TrimSpace(r.FormValue("file_type")),
			Description: strings.TrimSpace(r.FormValue("description")),
		}

		if err := h.store.UpdateAsset(asset); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/assets", http.StatusSeeOther)
		return

	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) DeleteAsset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id, err := strconv.ParseInt(r.URL.Query().Get("id"), 10, 64)
	if err != nil || id <= 0 {
		http.Error(w, "invalid asset id", http.StatusBadRequest)
		return
	}

	if err := h.store.DeleteAsset(id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/assets", http.StatusSeeOther)
}
