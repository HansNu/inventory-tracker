package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	"inventory-tracker/handler"
)

var db *pgx.Conn //global db connection

func main() {
	db, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer db.Close(context.Background()) //defer conn.close : db connection when main finishes running
	handler.DB = db

	r := gin.Default()
	r.POST("/addAsset", handler.AddAsset)
	r.GET("/getAssetList", handler.GetAssetList)
	r.PUT("/updateAsset", handler.UpdateAsset)
	r.DELETE("/deleteAsset", handler.DeleteAsset)

	r.Run(":8080")
}
