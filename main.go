package main

import (
	"context"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

var db *pgx.Conn //global db connection

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Println("Unable to connect to database:", err)
		return
	}
	defer conn.Close(context.Background()) //defer conn.close : db connection when main finishes running

	r := gin.Default()
	// r.POST("/assets", addAsset)
	r.Run(":8080")
}
