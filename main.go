package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"inventory-tracker/handler"
)

func main() {
	db, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL")) // change this
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer db.Close() // pgxpool.Close() doesn't take a context
	//defer conn.close : db connection when main finishes running
	h := &handler.Handler{DB: db}
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
	}))

	const api string = "api"
	r.POST(api+"/addAsset", h.AddAsset)
	r.POST(api+"/addAssetCategory", h.AddAssetCategory)
	r.POST(api+"/addCategoryGroup", h.AddCategoryGroup)

	r.GET(api+"/getAssetList", h.GetAssetList)
	r.GET(api+"/getAssetCategoryList", h.GetAssetCategoryList)
	r.GET(api+"/getAssetByAssetCode/:assetCode", h.GetAssetByAssetCode)
	r.GET(api+"/getCategoryGroup", h.GetCategoryGroup)

	r.PUT(api+"/updateAsset", h.UpdateAsset)

	r.DELETE(api+"/deleteAssetByAssetCode", h.DeleteAssetByAssetCode)
	r.DELETE(api+"/deleteAssetCategoryById", h.DeleteAssetCategoryById)
	r.DELETE(api+"/deleteCategoryGroup", h.DeleteCategoryGroupById)

	r.Run(":8080")
}
