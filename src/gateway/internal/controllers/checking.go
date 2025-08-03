package controllers

import (
	"context"

	"github.com/Fl0rencess720/Majula/src/gateway/internal/models"
	pb "github.com/Fl0rencess720/Majula/src/idl/checking"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type CheckingRepo interface {
}

type CheckingUsecase struct {
	repo           CheckingRepo
	checkingClient pb.FactCheckingClient
}

func NewCheckingUsecase(repo CheckingRepo, cc pb.FactCheckingClient) *CheckingUsecase {
	return &CheckingUsecase{repo: repo, checkingClient: cc}
}

func (u *CheckingUsecase) Check(c *gin.Context) {
	req := models.CheckingReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		zap.L().Error("ShouldBindJSON failed:", zap.Error(err))
		ErrorResponse(c, ServerError)
		return
	}
	ctx := context.Background()
	result, err := u.checkingClient.Check(ctx, &pb.CheckRequest{Text: req.Text})
	if err != nil {
		zap.L().Error("check failed:", zap.Error(err))
		ErrorResponse(c, ServerError)
		return
	}
	resp := []models.CheckingResp{}
	for _, e := range result.Results {
		resp = append(resp, models.CheckingResp{
			OriginalText: e.OriginalText,
			Sources:      e.Sources,
			Result:       e.Result,
			Reason:       e.Reason,
		})
	}
	SuccessResponse(c, resp)
}
