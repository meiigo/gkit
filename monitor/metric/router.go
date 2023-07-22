package metric

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Router(rg *gin.RouterGroup) {
	rg.GET("/metrics", gin.WrapH(promhttp.Handler()))
}
