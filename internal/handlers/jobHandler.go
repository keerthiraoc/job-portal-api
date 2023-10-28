package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/keerthiraoc/job-portal-api/internal/auth"
	"github.com/keerthiraoc/job-portal-api/internal/middleware"
	"github.com/keerthiraoc/job-portal-api/internal/models"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	"github.com/gin-gonic/gin"
)

func (h *handler) RegisterCompany(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceID missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	claims, ok := ctx.Value(auth.AuthKey).(jwt.RegisteredClaims)
	if !ok {
		log.Error().Str("Trace Id", traceID).Msg("user not logged in")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	var newComp models.NewCompany
	err := json.NewDecoder(c.Request.Body).Decode(&newComp)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	validate := validator.New()
	err = validate.Struct(newComp)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide company name, location"})
		return
	}

	uid, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.
			StatusInternalServerError)})
		return
	}

	cmp, err := h.s.CreateCompany(ctx, newComp, uint(uid))
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "company creation failed"})
		return
	}

	c.JSON(http.StatusOK, cmp)
}

func (h *handler) GetCompanies(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceID missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	cmp, err := h.s.FetchCompanies(ctx)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": "companies fetch failed"})
		return
	}

	c.JSON(http.StatusOK, cmp)
}

func (h *handler) GetCompany(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceID missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	stringCID := c.Param("companyID")
	cID, err := strconv.ParseUint(stringCID, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.
			StatusInternalServerError)})
		return
	}

	cmp, err := h.s.FetchCompany(ctx, uint(cID))
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cmp)
}

func (h *handler) AddJobToCompany(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceID missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	claims, ok := ctx.Value(auth.AuthKey).(jwt.RegisteredClaims)
	if !ok {
		log.Error().Str("Trace Id", traceID).Msg("user not logged in")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": http.StatusText(http.StatusUnauthorized)})
		return
	}

	var nj models.NewJob
	err := json.NewDecoder(c.Request.Body).Decode(&nj)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	validate := validator.New()
	err = validate.Struct(nj)
	if err != nil {
		log.Error().Err(err).Str("Trace ID", traceID)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"msg": "please provide title, salary and description"})
		return
	}

	uid, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.
			StatusInternalServerError)})
		return
	}

	stringCID := c.Param("companyID")
	cID, err := strconv.ParseUint(stringCID, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.
			StatusInternalServerError)})
		return
	}

	job, err := h.s.AddJobToCompany(ctx, nj, uint(uid), uint(cID))
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}

func (h *handler) ListCompanyJobs(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceID missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	stringCID := c.Param("companyID")
	cID, err := strconv.ParseUint(stringCID, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.
			StatusInternalServerError)})
		return
	}

	jobs, err := h.s.ListCompanyJobs(ctx, uint(cID))
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *handler) ListJobs(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceID missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	jobs, err := h.s.ListJobs(ctx)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jobs)
}

func (h *handler) GetJobByID(c *gin.Context) {
	ctx := c.Request.Context()
	traceID, ok := ctx.Value(middleware.TraceIDKey).(string)
	if !ok {
		log.Error().Msg("traceID missing from context")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": http.StatusText(http.StatusInternalServerError)})
		return
	}

	stringJID := c.Param("jobID")
	jobID, err := strconv.ParseUint(stringJID, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"msg": http.StatusText(http.
			StatusInternalServerError)})
		return
	}

	job, err := h.s.GetJobByID(ctx, uint(jobID))
	if err != nil {
		log.Error().Err(err).Str("Trace Id", traceID)
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"msg": err.Error()})
		return
	}
	c.JSON(http.StatusOK, job)
}
