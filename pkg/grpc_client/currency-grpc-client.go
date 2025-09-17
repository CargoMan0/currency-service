package grpc_client

import (
	"fmt"
	"github.com/BernsteinMondy/currency-service/pkg/currency"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func NewCurrencyGRPCClient(addr string) (currency.CurrencyServiceClient, *grpc.ClientConn, error) {
	dialOptions := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(addr, dialOptions...)
	if err != nil {
		return nil, nil, fmt.Errorf("grpc.NewClient: %w", err)
	}

	client := currency.NewCurrencyServiceClient(conn)

	return client, conn, nil
}
