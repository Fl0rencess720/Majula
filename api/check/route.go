package check

import (
	"github.com/Fl0rencess720/Majula/internal/controllers"
	"github.com/gin-gonic/gin"
)

func InitApi(group *gin.RouterGroup, cu *controllers.CheckingUsecase) {
	group.GET("/checking", cu.Check)
}
