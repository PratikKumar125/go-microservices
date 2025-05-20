package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PratikKumar125/go-microservices/pkg/logging"
	"github.com/PratikKumar125/go-microservices/users/internal/config"
	"github.com/PratikKumar125/go-microservices/users/internal/db"
	"github.com/PratikKumar125/go-microservices/users/internal/handler"
	"github.com/PratikKumar125/go-microservices/users/internal/repositories"
	"github.com/PratikKumar125/go-microservices/users/internal/service"
	usersgrpc "github.com/PratikKumar125/go-microservices/users/usersrpc"
	"github.com/knadh/koanf/v2"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// 0 Logger setup
	logger := logging.NewLogger("users", "debug")

	// 1 Koanf setup
	var k = koanf.New(".")
	appConfig, err := config.NewAppConfig(k, "/Users/pratikkumar/Downloads/Personal/golang/Microservices/users/dev.env.yaml")
	if err != nil {
		panic(err)
	}

	// 2. Database Init
	db, err := db.NewDatabase(ctx, appConfig, logger)
	if err != nil {
		logger.Error("Failed to initialize database", "error", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		logger.Error("Failed to ping database", "error", err)
		panic(err)
	}

	// 3. Repository Init
	userRepo := repositories.NewUserRepository(appConfig, db)

	// 4. User Service Init
	userService := service.NewUserService(appConfig, logger, userRepo)

	// 5. User Handler Init
	userHandler := handler.NewUserHandler(appConfig, userService)

	// 6. Initialize HTTP Server
	server := &http.Server{
		Addr:         appConfig.ConfigService.String("app.port"),
		Handler:      userHandler.InitRoutes(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// 7. Start server as gorutine
	go func() {
		logger.Info("Started user service")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed: ", err)
		}
	}()

	// 8. Start gRPC server
	grpcAddr := appConfig.ConfigService.String("app.grpc_port")
	lis, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		logger.Error("Failed to listen for gRPC", "error", err)
	}

	grpcServer := grpc.NewServer()
	grpcService := usersgrpc.NewUserGrpcServer(userService)
	usersgrpc.RegisterUserRpcServiceServer(grpcServer, grpcService)

	go func() {
		logger.Info("Started gRPC server")
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC server failed", "error", err)
		}
	}()

	// 9. Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "error", err)
	}
	logger.Info("Server Stopped")
}
