package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"inventory-tracker/handler"
)

func main() {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer db.Close(context.Background()) //defer conn.close : db connection when main finishes running
	h := &handler.Handler{DB: db}
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders: []string{"Content-Type"},
	}))

	r.POST("/addAsset", h.AddAsset)

	r.GET("/getAssetList", h.GetAssetList)
	r.GET("/getAssetCategoryList", h.GetAssetCategoryList)

	r.PUT("/updateAsset", h.UpdateAsset)

	r.DELETE("/deleteAsset", h.DeleteAsset)

	r.Run(":8080")
}
