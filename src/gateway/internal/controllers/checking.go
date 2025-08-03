package controllers

import (
	"context"

	"github.com/Fl0rencess720/Majula/src/gateway/internal/pkgs/agent"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CheckingRepo interface {
}

type CheckingUsecase struct {
	repo CheckingRepo
}

type CheckingReq struct {
	Text string `json:"text"`
}

func NewCheckingUsecase(repo CheckingRepo) *CheckingUsecase {
	return &CheckingUsecase{repo: repo}
}

func (u *CheckingUsecase) Check(c *gin.Context) {
	req := CheckingReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("ShouldBindJSON failed:", zap.Error(err))
		ErrorResponse(c, ServerError)
		return
	}
	ctx := context.Background()

	agent, err := agent.NewCheckingAgent(ctx)
	if err != nil {
		zap.L().Error("NewCheckingAgent failed:", zap.Error(err))
		ErrorResponse(c, ServerError)
		return
	}
	result, err := agent.Run(ctx, req.Text)
	if err != nil {
		zap.L().Error("CheckingAgent failed:", zap.Error(err))
		ErrorResponse(c, ServerError)
		return
	}
	SuccessResponse(c, result)
}
