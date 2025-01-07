package runner

import (
	"github.com/Le-BlitzZz/real-time-chat-app/config"
	"github.com/gin-gonic/gin"
)

func Run(router *gin.Engine, conf *config.Configuration) error {
	return router.Run(conf.Server.ListenAddr + ":" + conf.Server.Port)
}
