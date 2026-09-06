package main

import (
	"context"
	"fmt"
	"inventory-tracker/handler"
	"inventory-tracker/middleware"
	"inventory-tracker/repository"
	service "inventory-tracker/services"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-delve/delve/pkg/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer db.Close()

	config.LoadConfig()

	repo := repository.NewAssetRepo(db)
	svc := service.NewAssetService(repo)
	h := &handler.Handler{
		DB:      db,
		Service: svc,
	}
	r := gin.Default()
	r.SetTrustedProxies([]string{"127.0.0.1"})

	allowedOrigin := os.Getenv("FRONTEND_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:5173"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{allowedOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	const api string = "api"

	auth := r.Group(api + "/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}
	r.GET(api+"/getCategoryGroup", h.GetCategoryGroup)

	protected := r.Group(api)
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", h.GetCurrentUser)

		protected.POST("/addAsset", h.AddAsset)
		protected.POST("/addAssetCategory", h.AddAssetCategory)
		protected.POST("/addCategoryGroup", h.AddCategoryGroup)

		protected.GET("/getAssetList", h.GetAssetList)
		protected.GET("/getAssetCategoryList", h.GetAssetCategoryList)
		protected.GET("/getAssetByAssetCode/:assetCode", h.GetAssetByAssetCode)

		protected.PUT("/updateAsset", h.UpdateAsset)

		protected.DELETE("/deleteAssetByAssetCode", h.DeleteAssetByAssetCode)
		protected.DELETE("/deleteAssetCategoryById", h.DeleteAssetCategoryById)
	}

	admin := r.Group(api)
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("IT"))
	{
		admin.DELETE("/deleteCategoryGroup", h.DeleteCategoryGroupById)
	}

	r.Run(":8080")
}
