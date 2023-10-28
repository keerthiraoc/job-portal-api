package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type key string

const TraceIDKey key = "traceID"

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := uuid.NewString()
		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, TraceIDKey, traceID)
		req := c.Request.WithContext(ctx)
		c.Request = req

		log.Info().Str("Trace Id", traceID).Str("Method", c.Request.Method).
			Str("URL Path", c.Request.URL.Path).Msg("request started")

		defer log.Info().Str("Trace Id", traceID).Str("Method", c.Request.Method).
			Str("URL Path", c.Request.URL.Path).Msg("request processing completed")

		c.Next()
	}
}
