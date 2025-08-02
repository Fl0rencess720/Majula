package main

import (
	"context"

	"github.com/Fl0rencess720/Majula/api"
	"github.com/Fl0rencess720/Majula/api/mcp"
	"github.com/Fl0rencess720/Majula/internal/conf"
	"github.com/Fl0rencess720/Majula/internal/controllers"
	"github.com/Fl0rencess720/Majula/internal/data"
	"github.com/Fl0rencess720/Majula/internal/pkgs/logging"
	"github.com/Fl0rencess720/Majula/internal/pkgs/profiling"
	"github.com/Fl0rencess720/Majula/internal/pkgs/tracing"
	ccb "github.com/cloudwego/eino-ext/callbacks/cozeloop"
	"github.com/cloudwego/eino/callbacks"
	"github.com/coze-dev/cozeloop-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	Name = "Majula.Gateway"
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
		panic(err)
	}
	defer client.Close(ctx)
	handler := ccb.NewLoopHandler(client)
	callbacks.AppendGlobalHandlers(handler)

	e := newSrv()
	mcp.Init(e)
	e.Run(viper.GetString("server.http.addr"))
}

func newSrv() *gin.Engine {
	checkingRepo := data.NewCheckingRepo()
	checkingUsecase := controllers.NewCheckingUsecase(checkingRepo)

	e := api.Init(checkingUsecase)
	return e
}
