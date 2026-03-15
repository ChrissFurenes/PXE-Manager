package admin

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"PXE-Manager/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	store     *storage.Store
	uploadDir string
}

func NewHandler(store *storage.Store, uploadDir string) *Handler {
	return &Handler{
		store:     store,
		uploadDir: uploadDir,
	}
}

func (h *Handler) Dashboard(c *gin.Context) {
	profiles, _ := h.store.ListProfiles()
	clients, _ := h.store.ListClients()
	assets, _ := h.store.ListAssets()
	defaultProfileID, _ := h.store.GetDefaultProfileID()

	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"ProfilesCount":    len(profiles),
		"ClientsCount":     len(clients),
		"AssetsCount":      len(assets),
		"DefaultProfileID": defaultProfileID,
	})
}

func (h *Handler) Profiles(c *gin.Context) {
	profiles, err := h.store.ListProfiles()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	defaultProfileID, _ := h.store.GetDefaultProfileID()

	c.HTML(http.StatusOK, "profiles.html", gin.H{
		"Profiles":         profiles,
		"DefaultProfileID": defaultProfileID,
	})
}

func (h *Handler) NewProfileForm(c *gin.Context) {
	assets, _ := h.store.ListAssets()

	c.HTML(http.StatusOK, "profile_form.html", gin.H{
		"Assets":  assets,
		"Profile": nil,
		"IsEdit":  false,
	})
}

func (h *Handler) NewProfile(c *gin.Context) {
	enabled := c.PostForm("enabled") == "on"

	p := &storage.Profile{
		Name:      strings.TrimSpace(c.PostForm("name")),
		BootMode:  strings.TrimSpace(c.PostForm("boot_mode")),
		BootType:  strings.TrimSpace(c.PostForm("boot_type")),
		Kernel:    strings.TrimSpace(c.PostForm("kernel")),
		Initrd:    strings.TrimSpace(c.PostForm("initrd")),
		ImagePath: strings.TrimSpace(c.PostForm("image_path")),
		Cmdline:   strings.TrimSpace(c.PostForm("cmdline")),
		Enabled:   enabled,
	}

	if err := h.store.CreateProfile(p); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/profiles/")
}

func (h *Handler) EditProfileForm(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid profile id")
		return
	}

	profile, err := h.store.GetProfile(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	assets, _ := h.store.ListAssets()

	c.HTML(http.StatusOK, "profile_form.html", gin.H{
		"Assets":  assets,
		"Profile": profile,
		"IsEdit":  true,
	})
}

func (h *Handler) EditProfile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid profile id")
		return
	}

	enabled := c.PostForm("enabled") == "on"

	p := &storage.Profile{
		ID:        id,
		Name:      strings.TrimSpace(c.PostForm("name")),
		BootMode:  strings.TrimSpace(c.PostForm("boot_mode")),
		BootType:  strings.TrimSpace(c.PostForm("boot_type")),
		Kernel:    strings.TrimSpace(c.PostForm("kernel")),
		Initrd:    strings.TrimSpace(c.PostForm("initrd")),
		ImagePath: strings.TrimSpace(c.PostForm("image_path")),
		Cmdline:   strings.TrimSpace(c.PostForm("cmdline")),
		Enabled:   enabled,
	}

	if err := h.store.UpdateProfile(p); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/profiles/")
}

func (h *Handler) DeleteProfile(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid profile id")
		return
	}

	if err := h.store.DeleteProfile(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/profiles/")
}

func (h *Handler) Clients(c *gin.Context) {
	clients, err := h.store.ListClients()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "clients.html", gin.H{
		"Clients": clients,
	})
}

func (h *Handler) NewClientForm(c *gin.Context) {
	profiles, _ := h.store.ListProfiles()

	c.HTML(http.StatusOK, "client_form.html", gin.H{
		"Profiles": profiles,
		"Client":   nil,
		"IsEdit":   false,
	})
}

func (h *Handler) NewClient(c *gin.Context) {
	showMenu := c.PostForm("show_menu") == "on"
	profileID, _ := strconv.ParseInt(c.PostForm("profile_id"), 10, 64)

	client := &storage.Client{
		MAC:         strings.TrimSpace(strings.ToLower(c.PostForm("mac"))),
		Hostname:    strings.TrimSpace(c.PostForm("hostname")),
		ProfileID:   profileID,
		ShowMenu:    showMenu,
		Description: strings.TrimSpace(c.PostForm("description")),
	}

	if err := h.store.CreateClient(client); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/clients/")
}

func (h *Handler) EditClientForm(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid client id")
		return
	}

	client, err := h.store.GetClient(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	profiles, _ := h.store.ListProfiles()

	c.HTML(http.StatusOK, "client_form.html", gin.H{
		"Profiles": profiles,
		"Client":   client,
		"IsEdit":   true,
	})
}

func (h *Handler) EditClient(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid client id")
		return
	}

	showMenu := c.PostForm("show_menu") == "on"
	profileID, _ := strconv.ParseInt(c.PostForm("profile_id"), 10, 64)

	client := &storage.Client{
		ID:          id,
		MAC:         strings.TrimSpace(strings.ToLower(c.PostForm("mac"))),
		Hostname:    strings.TrimSpace(c.PostForm("hostname")),
		ProfileID:   profileID,
		ShowMenu:    showMenu,
		Description: strings.TrimSpace(c.PostForm("description")),
	}

	if err := h.store.UpdateClient(client); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/clients/")
}

func (h *Handler) DeleteClient(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid client id")
		return
	}

	if err := h.store.DeleteClient(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/clients/")
}

func (h *Handler) Assets(c *gin.Context) {
	assets, err := h.store.ListAssets()
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.HTML(http.StatusOK, "assets.html", gin.H{
		"Assets": assets,
	})
}

func (h *Handler) NewAssetForm(c *gin.Context) {
	c.HTML(http.StatusOK, "asset_form.html", gin.H{
		"Asset":  nil,
		"IsEdit": false,
	})
}

func (h *Handler) NewAsset(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	description := strings.TrimSpace(c.PostForm("description"))
	fileType := strings.TrimSpace(c.PostForm("file_type"))

	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	if name == "" {
		name = file.Filename
	}

	if err := os.MkdirAll(h.uploadDir, 0o755); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	safeName := filepath.Base(file.Filename)
	dstPath := filepath.Join(h.uploadDir, safeName)

	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	asset := &storage.Asset{
		Name:        name,
		FileName:    safeName,
		FilePath:    "/images/" + safeName,
		FileType:    fileType,
		SizeBytes:   file.Size,
		Description: description,
	}

	if err := h.store.CreateAsset(asset); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/assets/")
}

func (h *Handler) EditAssetForm(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid asset id")
		return
	}

	asset, err := h.store.GetAsset(id)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	c.HTML(http.StatusOK, "asset_form.html", gin.H{
		"Asset":  asset,
		"IsEdit": true,
	})
}

func (h *Handler) EditAsset(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid asset id")
		return
	}

	asset := &storage.Asset{
		ID:          id,
		Name:        strings.TrimSpace(c.PostForm("name")),
		FileType:    strings.TrimSpace(c.PostForm("file_type")),
		Description: strings.TrimSpace(c.PostForm("description")),
	}

	if err := h.store.UpdateAsset(asset); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/assets/")
}

func (h *Handler) DeleteAsset(c *gin.Context) {
	id, err := strconv.ParseInt(c.Query("id"), 10, 64)
	if err != nil || id <= 0 {
		c.String(http.StatusBadRequest, "invalid asset id")
		return
	}

	if err := h.store.DeleteAsset(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/assets/")
}

func (h *Handler) SetDefaultProfile(c *gin.Context) {
	id, err := strconv.ParseInt(c.PostForm("profile_id"), 10, 64)
	if err != nil || id < 0 {
		c.String(http.StatusBadRequest, "invalid profile id")
		return
	}

	if err := h.store.SetDefaultProfileID(id); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusSeeOther, "/admin/profiles/")
}
