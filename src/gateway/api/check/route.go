package check

import (
	"github.com/Fl0rencess720/Majula/src/gateway/internal/controllers"
	"github.com/gin-gonic/gin"
)

func InitApi(group *gin.RouterGroup, cu *controllers.CheckingUsecase) {
	group.POST("/checking", cu.Check)
}
