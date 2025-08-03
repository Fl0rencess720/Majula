package api

import (
	"time"

	"github.com/Fl0rencess720/Majula/src/gateway/api/check"
	"github.com/Fl0rencess720/Majula/src/gateway/internal/controllers"
	"github.com/Fl0rencess720/Majula/src/gateway/internal/middlewares"
	ginZap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Init(cu *controllers.CheckingUsecase) *gin.Engine {
	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery(), ginZap.Ginzap(zap.L(), time.RFC3339, false), ginZap.RecoveryWithZap(zap.L(), false))

	app := e.Group("/api", middlewares.Cors())
	{
		check.InitApi(app, cu)
	}
	return e
}
