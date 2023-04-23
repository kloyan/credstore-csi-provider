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
	"github.com/kloyan/credstore-csi-provider/internal/version"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var Logger *zap.SugaredLogger = initLogger()

func initLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	return sugar
}

func main() {
	var serviceKeyPath, providerPath string

	flag.StringVar(&serviceKeyPath, "service-key-path", "/tmp/service-key.json", "Path to file which contains the service key")
	flag.StringVar(&providerPath, "provider-path", "/tmp", "Path to directory in which the provider unix domain socket shall be created")
	flag.Parse()

	ver := version.GetVersion()

	Logger.Infow("initializing credstore provider",
		"version", ver.BuildVersion,
		"commit", ver.GitCommit,
	)

	if err := startServer(serviceKeyPath, providerPath); err != nil {
		Logger.Errorw("error running grpc server", "err", err)
		os.Exit(1)
	}
}

func startServer(serviceKeyPath, providerPath string) error {
	serviceKey, err := readServiceKey(serviceKeyPath)
	if err != nil {
		return err
	}

	encryptor, err := client.NewJWEDecryptor(serviceKey)
	if err != nil {
		return err
	}

	client, err := client.NewClient(serviceKey, encryptor, 3*time.Second)
	if err != nil {
		return err
	}

	provider := provider.NewProvider(client)

	interceptor := grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		Logger.Infow("processing grpc request", "grpc.method", info.FullMethod)
		resp, err := handler(ctx, req)
		Logger.Infow("finished grpc request",
			"grpc.method", info.FullMethod,
			"grpc.code", status.Code(err),
			"err", err,
		)

		return resp, err
	})

	server := server.NewServer(provider, providerPath, interceptor)

	c := make(chan os.Signal, 2)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-c
		Logger.Infof("caught os signal %s, shutting down", sig)
		server.Stop()
	}()

	Logger.Info("starting grpc server")
	if err := server.Start(); err != nil {
		return err
	}

	return nil
}

func readServiceKey(serviceKeyPath string) (config.ServiceKey, error) {
	jsonBytes, err := os.ReadFile(serviceKeyPath)
	if err != nil {
		return config.ServiceKey{}, err
	}

	return config.ParseServiceKey(jsonBytes)
}
