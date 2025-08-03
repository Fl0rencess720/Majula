package main

import (
	"context"

	"github.com/Fl0rencess720/Majula/src/common/conf"
	"github.com/Fl0rencess720/Majula/src/common/logging"
	"github.com/Fl0rencess720/Majula/src/common/profiling"
	"github.com/Fl0rencess720/Majula/src/common/tracing"
	"github.com/Fl0rencess720/Majula/src/services/checking/internal/service"

	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/coze-dev/cozeloop-go"
	"go.uber.org/zap"
)

var (
	Name = "Majula.Service.Checking"
)

func init() {
	conf.Init()

	logging.Init()

	profiling.InitPyroscope(Name)

}

func main() {
	ctx := context.Background()
	tp, err := tracing.SetTraceProvider(Name)
	if err != nil {
		zap.L().Panic("tracing init err: %s", zap.Error(err))
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			zap.L().Error("trace provider shut down err: %s", zap.Error(err))
		}
	}()

	client, err := cozeloop.NewClient()
	if err != nil {
		zap.L().Panic("cozeloop init err: %s", zap.Error(err))
	}
	defer client.Close(ctx)
	handler := ccb.NewLoopHandler(client)
	callbacks.AppendGlobalHandlers(handler)

	grpcService, err := service.NewCheckingService(Name)
	if err != nil {
		zap.L().Panic("Failed to create gRPC service", zap.Error(err))
	}

	if err := grpcService.Start(); err != nil {
		zap.L().Panic("Failed to start service", zap.Error(err))
	}

	grpcService.WaitForShutdown()

	if err := grpcService.Stop(); err != nil {
		zap.L().Error("Error stopping service", zap.Error(err))
	}

	zap.L().Info("Server exit")
}
