package main

import (
	"SocialMediaGO/Controllers"
	"SocialMediaGO/middleware"
	"SocialMediaGO/migrate"
	"github.com/gin-gonic/gin"
)

func init() {
	migrate.Migrate()
}

func main() {
	r := gin.Default()

	r.POST("/register", Controllers.Register)

	r.POST("/login", Controllers.Login)

	r.POST("/logout", middleware.RequireAuth, Controllers.Logout)

	r.PUT("/update/", middleware.RequireAuth, Controllers.UpdateProfile)

	r.GET("/me", middleware.RequireAuth, Controllers.Me)

	r.DELETE("/deactivate", middleware.RequireAuth, Controllers.Deactivate)

	r.POST("post", middleware.RequireAuth, Controllers.Post)

	r.DELETE("/delete_post/:post_id", middleware.RequireAuth, Controllers.DeletePost)

	r.GET("/get_post/:post_id", middleware.RequireAuth, Controllers.GetPost)

	r.GET("/posts", middleware.RequireAuth, Controllers.GetMyPostsList)

	r.Run()
}
