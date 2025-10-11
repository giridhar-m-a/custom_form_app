package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/giridhar-m-a/custom_form_app/internal/api"
	"github.com/giridhar-m-a/custom_form_app/internal/cache"
	"github.com/giridhar-m-a/custom_form_app/internal/db"
	"github.com/giridhar-m-a/custom_form_app/internal/services"
	"github.com/giridhar-m-a/custom_form_app/internal/utils"

	_ "github.com/giridhar-m-a/custom_form_app/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title My API
// @version 1.0
// @description Backend API for Custom Form Application
// @contact.name API Support
// @contact.url http://giridhar.dev/
// @contact.email m.a.giridhar08@gmail.com
// @host localhost:8000
// @BasePath /api/v1
// @schemes https http
// @securityDefinitions.apikey BearerAuth
// @description JWT Authorization header using the Bearer scheme.
// @in header
// @name Authorization
func main() {
	// Get configuration from environment
	port := utils.GetEnv("APP_PORT", "8080")
	allowedOrigins := []string{
		utils.GetEnv("FRONTEND_URL", "http://localhost:3000"),
		"http://localhost:8000",
	}

	// Initialize database
	log.Println("Initializing database...")
	db.InitDB()

	// Initialize cache
	log.Println("Initializing cache...")
	cache.Init()
	defer cache.Close()

	// Initialize MinIO
	log.Println("Initializing MinIO...")
	services.InitMinio()

	// Set Gin mode based on environment
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	r := gin.Default()

	// Configure CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Server is running",
		})
	})

	// Swagger endpoints
	r.GET("/api/docs", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/api/docs/index.html")
	})
	r.GET("/api/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Register API routes
	api.RegisterRoutes(r)

	// Create HTTP server
	srv := &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting server on port %s...", port)
		log.Printf("Swagger documentation available at: http://localhost:%s/api/docs", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// Kill (no param) default sends syscall.SIGTERM
	// Kill -2 is syscall.SIGINT (Ctrl+C)
	// Kill -9 is syscall.SIGKILL (cannot be caught)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}

