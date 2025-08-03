package biz

import (
	"context"

	"github.com/Fl0rencess720/Majula/src/services/checking/internal/models"
	"github.com/Fl0rencess720/Majula/src/services/checking/internal/pkgs/agent"
)

type CheckingRepo interface {
}

type CheckingUsecase struct {
	repo CheckingRepo
}

func NewCheckingUsecase(repo CheckingRepo) *CheckingUsecase {
	return &CheckingUsecase{repo: repo}
}

func (u *CheckingUsecase) Check(text string) ([]models.CheckingResult, error) {
	ctx := context.Background()

	agent, err := agent.NewCheckingAgent(ctx)
	if err != nil {
		return nil, err
	}
	result, err := agent.Run(ctx, text)
	if err != nil {

		return nil, err
	}
	return result, nil

}
