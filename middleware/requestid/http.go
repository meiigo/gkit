package requestid

import (
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

var randSource = rand.New(rand.NewSource(time.Now().UnixNano()))

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestID string
		// 1. from context
		requestID, _ = FromContext(c.Request.Context())
		if requestID != "" {
			c.Writer.Header().Set(key, requestID)
			c.Next()
			return
		}

		// 2. from incoming
		requestID = c.GetHeader(key)
		if requestID != "" {
			c.Set(key, requestID)
			c.Next()
			return
		}

		// 3. generate
		var tid [16]byte
		randSource.Read(tid[:])
		requestID = hex.EncodeToString(tid[:])
		c.Set(key, requestID)
		c.Writer.Header().Set(key, requestID)
		c.Next()
	}
}
