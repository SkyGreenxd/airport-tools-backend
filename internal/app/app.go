package app

import (
	"airport-tools-backend/internal/config"
	v1 "airport-tools-backend/internal/delivery/v1"
	"airport-tools-backend/internal/infrastructure"
	"airport-tools-backend/internal/repository/postgres"
	"airport-tools-backend/internal/repository/yandex_s3"
	"airport-tools-backend/internal/server"
	"airport-tools-backend/internal/usecase"
	"airport-tools-backend/pkg/logger"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func Run() {
	pg, err := postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}

	if err := pg.Ping(); err != nil {
		log.Fatal(err)
	}

	if err := pg.RunMigrations(); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := pg.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	userRepo := postgres.NewUserRepository(pg.Db)
	cvScanDetailRepo := postgres.NewCvScanDetailRepository(pg.Db)
	cvScanRepo := postgres.NewCvScanRepository(pg.Db)
	toolSetRepo := postgres.NewToolSetRepository(pg.Db)
	toolTypeRepo := postgres.NewToolTypeRepository(pg.Db)
	transactionRepo := postgres.NewTransactionRepository(pg.Db)

	bucketName := os.Getenv("BUCKET_NAME")
	s3, err := yandex_s3.InitS3(bucketName)
	if err != nil {
		log.Fatal(err)
	}

	mlUrl := os.Getenv("ML_SERVICE_URL")
	imageStorage := infrastructure.NewImageStorage(s3)
	ml := infrastructure.NewMlGateway(http.DefaultClient, mlUrl, imageStorage)

	strConfidence := os.Getenv("CONFIDENCE")
	confidence, err := strconv.ParseFloat(strConfidence, 32)
	if err != nil {
		log.Fatal(err)
	}
	strCosineSim := os.Getenv("COSINE_SIM")
	cosineSim, err := strconv.ParseFloat(strCosineSim, 32)
	if err != nil {
		log.Fatal(err)
	}

	trRepo := postgres.NewTransactionResolutionsRepo(pg.Db)
	loger := logger.NewSlogLogger()
	roleRepo := postgres.NewRoleRepo(pg.Db)
	service := usecase.NewService(userRepo, cvScanRepo, cvScanDetailRepo, toolTypeRepo, transactionRepo, ml, imageStorage, toolSetRepo, float32(confidence), float32(cosineSim), trRepo, loger, roleRepo)

	handler := v1.NewHandler(service)

	r := gin.Default()
	api := r.Group("/api")
	handler.Init(api)

	serverConfig := config.LoadHttpServerConfig()
	server := server.NewServer(r, serverConfig)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		log.Printf("starting server on port %s", serverConfig.Port)
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped gracefully")
}
