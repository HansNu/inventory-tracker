package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"inventory-tracker/handler"
	"inventory-tracker/middleware"
	"inventory-tracker/repository"
	service "inventory-tracker/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	_ "inventory-tracker/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title       Inventory Tracker API
// @version     1.0
// @description IT asset management API. Most routes require a bearer token from /auth/login.
// @BasePath    /api
// @securityDefinitions.apikey BearerAuth
// @in          header
// @name        Authorization
// @description Type "Bearer " followed by the token returned by /auth/login.

func main() {
	// .env is a convenience for local runs only. In production the variables
	// come from the platform, so a missing file here is expected, not an error.
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	db, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("DATABASE_URL is not a valid connection string: %v", err)
	}
	defer db.Close()

	// pgxpool.New only parses — it does not dial. Without this ping a bad
	// connection string surfaces as a 500 on the first request instead of a
	// refusal to start.
	if err := db.Ping(context.Background()); err != nil {
		log.Fatalf("cannot reach database at startup: %v", err)
	}

	// Composition root: the one place that knows which concrete implementation
	// sits behind each interface. Everything below here depends on interfaces.
	assetRepo := repository.NewAssetRepo(db)
	userRepo := repository.NewUserRepo(db)
	categoryRepo := repository.NewCategoryRepo(db)
	catGroupRepo := repository.NewCategoryGroupRepo(db)

	h := &handler.Handler{
		Service:        service.NewAssetService(assetRepo),
		Auth:           service.NewAuthService(userRepo),
		Categories:     categoryRepo,
		CategoryGroups: catGroupRepo,
	}

	r := gin.Default()
	r.SetTrustedProxies(nil)

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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("/health", func(c *gin.Context) {
		if err := db.Ping(c.Request.Context()); err != nil {
			log.Printf("health: db ping failed: %v", err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "degraded"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	auth := r.Group(api + "/auth")
	{
		auth.POST("/register", h.Register)
		auth.POST("/login", h.Login)
	}

	// TODO: this is the only read outside AuthMiddleware. The frontend needs it
	// before login to populate a dropdown on the register form; if that stops
	// being true, move it into the protected group.
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

		protected.PUT("/updateAsset/:id", h.UpdateAsset)

		protected.DELETE("/deleteAssetByAssetCode/:assetCode", h.DeleteAssetByAssetCode)
		protected.DELETE("/deleteAssetCategoryById/:id", h.DeleteAssetCategoryById)
	}

	admin := r.Group(api)
	admin.Use(middleware.AuthMiddleware(), middleware.RoleMiddleware("IT"))
	{
		admin.DELETE("/deleteCategoryGroup/:id", h.DeleteCategoryGroupById)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
