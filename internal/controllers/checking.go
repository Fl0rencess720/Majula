package controllers

import "github.com/gin-gonic/gin"

type CheckingUsecase struct {
}

func NewCheckingUsecase() *CheckingUsecase {
	return &CheckingUsecase{}
}

func (u *CheckingUsecase) Check(c *gin.Context) {

}
