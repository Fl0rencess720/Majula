package main

import (
	"context"

	"github.com/Fl0rencess720/Majula/src/common/conf"
	"github.com/Fl0rencess720/Majula/src/common/logging"
	"github.com/Fl0rencess720/Majula/src/common/profiling"
	"github.com/Fl0rencess720/Majula/src/common/registry"
	"github.com/Fl0rencess720/Majula/src/common/tracing"
	"github.com/Fl0rencess720/Majula/src/gateway/internal/controllers"
	"github.com/Fl0rencess720/Majula/src/gateway/internal/data"
	api "github.com/Fl0rencess720/Majula/src/gateway/service"
	"github.com/Fl0rencess720/Majula/src/gateway/service/mcp"

	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/coze-dev/cozeloop-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Name = "Majula.Gateway"
	ID   = ""
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
		zap.L().Panic("cozeloop client creation err: %s", zap.Error(err))
	}
	defer client.Close(ctx)
	handler := ccb.NewLoopHandler(client)
	callbacks.AppendGlobalHandlers(handler)

	e, err := newSrv()
	if err != nil {
		zap.L().Panic("server startup err: %s", zap.Error(err))
	}
	mcp.Init(e)
	e.Run(viper.GetString("server.http.addr"))
}

func newSrv() (*gin.Engine, error) {
	consulClient, err := registry.NewConsulClient(viper.GetString("CONSUL_ADDR"))
	if err != nil {
		return nil, err
	}

	cc, err := data.NewCheckingClient(viper.GetString("service.checking"))
	if err != nil {
		return nil, err
	}
	checkingRepo := data.NewCheckingRepo(cc)
	checkingUsecase := controllers.NewCheckingUsecase(checkingRepo)

	serviceID, err := consulClient.RegisterService(Name)
	if err != nil {
		return nil, err
	}
	ID = serviceID
	consulClient.SetTTLHealthCheck()

	e := api.Init(checkingUsecase)
	return e, nil
}
