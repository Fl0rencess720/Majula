package service

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Fl0rencess720/Majula/src/common/registry"
	pb "github.com/Fl0rencess720/Majula/src/idl/checking"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type CheckingService struct {
	pb.UnimplementedFactCheckingServer
	serviceName string
	serviceID   string
	registry    *registry.ConsulClient
	server      *grpc.Server
	listener    net.Listener
}

func NewCheckingService(serviceName string) (*CheckingService, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", viper.GetInt("server.grpc.port")))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}
	server := grpc.NewServer()

	pb.RegisterFactCheckingServer(server, &CheckingService{})
	registry, err := registry.NewConsulClient(viper.GetString("CONSUL_ADDR"))
	if err != nil {
		return nil, fmt.Errorf("failed to create registry: %w", err)
	}
	return &CheckingService{serviceName: serviceName, registry: registry, server: server, listener: lis}, nil
}

func (s *CheckingService) Start() error {
	serviceID, err := s.registry.RegisterService(s.serviceName)
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}
	s.serviceID = serviceID

	go s.registry.SetTTLHealthCheck()

	go func() {
		if err := s.server.Serve(s.listener); err != nil {
			zap.L().Error("Failed to serve", zap.Error(err))
		}
	}()
	return nil
}

func (s *CheckingService) Stop() error {
	if s.serviceID != "" {
		if err := s.registry.DeregisterService(s.serviceID); err != nil {
			zap.L().Error("Failed to deregister service",
				zap.String("service_id", s.serviceID),
				zap.Error(err))
		}
	}
	zap.L().Info("Shutting down gRPC server...")
	s.server.GracefulStop()
	return nil
}

func (s *CheckingService) WaitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Service is shutting down...")
}
