package main

import (
	"context"
	"fmt"
	"inventory-tracker/handler"
	"inventory-tracker/middleware"
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

	h := &handler.Handler{DB: db}
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
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
		protected.GET("/getCategoryGroup", h.GetCategoryGroup)

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
