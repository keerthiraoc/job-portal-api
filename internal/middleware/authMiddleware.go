package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/keerthiraoc/job-portal-api/internal/auth"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Mid struct {
	a *auth.Auth
}

func NewMid(a *auth.Auth) (*Mid, error) {
	if a == nil {
		return nil, errors.New("auth can't be nil")
	}
	return &Mid{a: a}, nil
}

func (m *Mid) Auth(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		traceID, ok := ctx.Value(TraceIDKey).(string)
		if !ok {
			log.Error().Msg("trace Id not present in the context")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
			return
		}

		authHeader := c.Request.Header.Get("Authorization")
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			err := errors.New("expected authorization header format: Bearer <token>")
			log.Error().Err(err).Str("Trace ID", traceID).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		claims, err := m.a.ValidateToken(parts[1])
		if err != nil {
			log.Error().Err(err).Str("Trace ID", traceID).Send()
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ctx = context.WithValue(ctx, auth.AuthKey, claims)
		req := c.Request.WithContext(ctx)
		c.Request = req

		next(c)
	}
}
