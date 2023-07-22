package env

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func Router(rg *gin.RouterGroup) {
	rg.GET("/env", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, os.Environ())
	})
}
