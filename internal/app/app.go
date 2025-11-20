package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilindan-dev/CommentTree/internal/config"
	"github.com/ilindan-dev/CommentTree/internal/storage/postgres"
	v1 "github.com/ilindan-dev/CommentTree/internal/transport/http/v1"
	"github.com/ilindan-dev/CommentTree/internal/usecase"
	"github.com/wb-go/wbf/zlog"
)

func Run() {
	zlog.InitConsole()
	logger := zlog.Logger
	logger.Info().Msg("App starting...")

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to load config")
	}
	logger.Info().Msg("Config loaded")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgPool, err := postgres.NewClient(ctx, cfg.DB)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer pgPool.Close()
	logger.Info().Msg("Database connected")

	repo := postgres.NewCommentRepository(pgPool, &logger)

	service := usecase.NewCommentService(repo, &logger)

	handler := v1.NewHandler(service, &logger)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	router.Static("/static", "./web/static")
	router.StaticFile("/", "./web/index.html")

	handler.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: router,
	}

	go func() {
		logger.Info().Msgf("Server starting on port %s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal().Err(err).Msg("Server startup failed")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Server shutting down...")

	ctxShutDown, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(ctxShutDown); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited properly")
}
