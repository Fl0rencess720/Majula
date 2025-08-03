package mcp

import (
	server "github.com/ckanthony/gin-mcp"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Init(e *gin.Engine) {
	mcp := server.New(e, &server.Config{
		Name:        "Fact Checking API",
		Description: "Input a text, get the fact checking results",
		BaseURL:     "http://localhost" + viper.GetString("server.http.addr"),
	})
	mcp.Mount("/mcp")
}
