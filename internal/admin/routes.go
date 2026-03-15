package admin

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine, h *Handler) {

	admin := router.Group("/admin")
	{

		profiles := admin.Group("/profiles")
		{
			profiles.GET("/", h.Profiles)
			profiles.GET("/new", h.NewProfileForm)
			profiles.POST("/new", h.NewProfile)
			profiles.GET("/edit", h.EditProfileForm)
			profiles.POST("/edit", h.EditProfile)
			profiles.POST("/delete", h.DeleteProfile)
		}

		clients := admin.Group("/clients")
		{
			clients.GET("/", h.Clients)
			clients.GET("/new", h.NewClientForm)
			clients.POST("/new", h.NewClient)
			clients.GET("/edit", h.EditClientForm)
			clients.POST("/edit", h.EditClient)
			clients.POST("/delete", h.DeleteClient)
		}

		assets := admin.Group("/assets")
		{
			assets.GET("/", h.Assets)
			assets.GET("/new", h.NewAssetForm)
			assets.POST("/new", h.NewAsset)
			assets.GET("/edit", h.EditAssetForm)
			assets.POST("/edit", h.EditAsset)
			assets.POST("/delete", h.DeleteAsset)
		}
		router.GET("/", h.Dashboard)
	}

}
