package app

import (
	auth "github.com/Shuhrat55/auth/pkg/api/g_rpc"
	"github.com/Shuhrat55/auth/pkg/logger"
	"github.com/Shuhrat55/auth/pkg/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
)

func StartGRPCServer() {
	port := os.Getenv("AUTH_GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logger.Logger.Fatal("Ошибка создания подключения для gRPC", zap.Error(err))
	}

	s := grpc.NewServer()
	auth.RegisterAuthServiceServer(s, &server.AuthServer{})

	logger.Logger.Debug("gRPC сервер стартует", zap.String("port", port))
	if err := s.Serve(lis); err != nil {
		logger.Logger.Error("Ошибка запуска gRPC сервера", zap.Error(err))
	}
}
