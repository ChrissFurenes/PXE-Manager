package admin

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/", h.Dashboard)

	mux.HandleFunc("/admin/profiles", h.Profiles)
	mux.HandleFunc("/admin/profiles/new", h.NewProfile)
	mux.HandleFunc("/admin/profiles/edit", h.EditProfile)
	mux.HandleFunc("/admin/profiles/delete", h.DeleteProfile)

	mux.HandleFunc("/admin/clients", h.Clients)
	mux.HandleFunc("/admin/clients/new", h.NewClient)
	mux.HandleFunc("/admin/clients/edit", h.EditClient)
	mux.HandleFunc("/admin/clients/delete", h.DeleteClient)

	mux.HandleFunc("/admin/assets", h.Assets)
	mux.HandleFunc("/admin/assets/new", h.NewAsset)
	mux.HandleFunc("/admin/assets/edit", h.EditAsset)
	mux.HandleFunc("/admin/assets/delete", h.DeleteAsset)

	mux.HandleFunc("/admin/settings/default-profile", h.SetDefaultProfile)
}
