package app

import (
	"log"

	"github.com/MXkodo/cash-server/config"
	"github.com/MXkodo/cash-server/internal/handlers"
	"github.com/MXkodo/cash-server/internal/middleware"
	"github.com/MXkodo/cash-server/internal/repo"
	"github.com/MXkodo/cash-server/internal/repo/postgresql"
	"github.com/MXkodo/cash-server/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func Run(cfg *config.Config, rdb *redis.Client) error {
	connstr := cfg.GetDBConnString()
	db, err := postgresql.NewDb(connstr)
	if err != nil {
		return err
	}
	defer postgresql.CloseDb(db)

	userRepo := repo.NewUserRepo(db)
	docRepo := repo.NewDocRepo(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret, rdb)
	docService := service.NewDocService(docRepo, rdb)

	authHandler := handlers.NewAuthHandler(authService, cfg.AdminToken)
	docHandler := handlers.NewDocHandler(docService)

	jwtMiddleware := middleware.NewJWTMiddleware(cfg.JWTSecret)

	router := gin.Default()

	authGroup := router.Group("/api")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/auth", authHandler.Authenticate)
		authGroup.DELETE("/auth/:token", authHandler.Logout)
	}

	docsGroup := router.Group("/api/docs").Use(jwtMiddleware.UserIDFromTokenMiddleware())
	{
		docsGroup.POST("", docHandler.UploadDocument)
		docsGroup.GET("", docHandler.GetDocuments)
		docsGroup.GET("/:id", docHandler.GetDocument)
		docsGroup.DELETE("/:id", docHandler.DeleteDocument)
	}


	if err := router.Run(cfg.ServerAddr); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
	return nil
}
