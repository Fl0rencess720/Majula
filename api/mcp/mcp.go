package mcp

import (
	server "github.com/ckanthony/gin-mcp"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Init(e *gin.Engine) {
	mcp := server.New(e, &server.Config{
		Name:        "My Simple API",
		Description: "An example API automatically exposed via MCP.",
		BaseURL:     "http://localhost" + viper.GetString("server.http.addr"),
	})
	mcp.Mount("/mcp")
}
