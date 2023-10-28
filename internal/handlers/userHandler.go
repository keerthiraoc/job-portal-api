package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/keerthiraoc/job-portal-api/internal/middleware"
	"github.com/keerthiraoc/job-portal-api/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
)

func (h *handler) Register(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceID missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	var nu models.NewUser

	err := json.NewDecoder(c.Request.Body).Decode(&nu)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	validate := validator.New()
	err = validate.Struct(nu)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide Name, Email and Password"})
		return
	}

	user, err := h.s.CreateUser(nu)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID).Msg("user registration error")
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "user registration failed"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *handler) Login(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceID missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	var lu models.UserLogin
	err := json.NewDecoder(c.Request.Body).Decode(&lu)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	validate := validator.New()
	err = validate.Struct(lu)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide Name, Email and Password"})
		return
	}

	claims, err := h.s.Authenticate(lu)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	var token models.NewToken
	token.Token, err = h.a.GenerateToken(claims)
	if err != nil {
		log.Error().Err(err).Msg("generating token")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.StatusInternalServerError)})
		return
	}

	c.JSON(http.StatusOK, token)
}
