package data

import (
	"fmt"
	"time"

	pb "github.com/Fl0rencess720/Majula/src/idl/checking"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type checkingRepo struct {
	checkingClient pb.FactCheckingClient
}

func NewCheckingRepo(cc pb.FactCheckingClient) *checkingRepo {
	return &checkingRepo{checkingClient: cc}
}

func NewCheckingClient(serviceName string) (pb.FactCheckingClient, error) {
	kacp := keepalive.ClientParameters{
		Time:                10 * time.Second,
		Timeout:             time.Second,
		PermitWithoutStream: false,
	}
	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s/%s?wait=15s", viper.GetString("CONSUL_ADDR"), serviceName),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		grpc.WithKeepaliveParams(kacp),
	)
	if err != nil {
		return nil, err
	}
	client := pb.NewFactCheckingClient(conn)
	return client, nil
}
