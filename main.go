package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kloyan/credstore-csi-provider/internal/client"
	"github.com/kloyan/credstore-csi-provider/internal/config"
	"github.com/kloyan/credstore-csi-provider/internal/provider"
	"github.com/kloyan/credstore-csi-provider/internal/server"
	"github.com/kloyan/credstore-csi-provider/internal/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "Enable debug logging")
	utils.InitLogger(debug)

	var serviceKeyPath string
	flag.StringVar(&serviceKeyPath, "service-key", "/tmp/service-key.json", "Path to file which contains the service key")

	var socketPath string
	flag.StringVar(&socketPath, "socket", "/tmp/credstore.sock", "Path to socket on which to listen for driver gRPC calls")

	if err := startServer(serviceKeyPath, socketPath); err != nil {
		utils.Logger.Errorw("error running grpc server", "err", err)
		os.Exit(1)
	}
}

func readServiceKey(serviceKeyPath string) (config.ServiceKey, error) {
	jsonBytes, err := os.ReadFile(serviceKeyPath)
	if err != nil {
		return config.ServiceKey{}, err
	}
	return config.ParseServiceKey(jsonBytes)
}

func startServer(serviceKeyPath, socketPath string) error {
	serviceKey, err := readServiceKey(serviceKeyPath)
	if err != nil {
		return err
	}

	client, err := client.NewClient(serviceKey, 3*time.Second)
	if err != nil {
		return err
	}

	provider := provider.NewProvider(client)

	interceptor := grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		utils.Logger.Infow("processing unary gRPC call", "grpc.method", info.FullMethod)
		resp, err := handler(ctx, req)
		utils.Logger.Infow("Finished unary gRPC call", "grpc.method", info.FullMethod, "grpc.code", status.Code(err), "err", err)
		return resp, err
	})

	server := server.NewServer(provider, socketPath, interceptor)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-c
		utils.Logger.Infof("caugh os signal %s, shutting down", sig)
		server.Stop()
	}()

	if err := server.Start(); err != nil {
		return err
	}
	return nil
}
