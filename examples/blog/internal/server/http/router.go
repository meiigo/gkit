package http

import (
	"github.com/gin-gonic/gin"
	"github.com/meiigo/gkit/middleware/requestid"
)

func (s *Server) Set(r *gin.Engine) {
	r.GET("_ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Pong"})
	})

	rg := r.Group("/article")
	rg.Use(
		requestid.GinMiddleware(),
		//timeoutMiddleware(200*time.Millisecond),
	)

	rg.GET("/detail", func(c *gin.Context) {
		c.JSON(200, "hello")
	})
}
