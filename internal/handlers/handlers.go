package handlers

import (
	"job-portal-api/internal/auth"
	"job-portal-api/internal/middleware"
	"job-portal-api/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type handler struct {
	s *services.Conn
	a *auth.Auth
}

func API(ms *services.Conn, a *auth.Auth) *gin.Engine {
	r := gin.New()

	h := handler{
		s: ms,
		a: a,
	}

	m, err := middleware.NewMid(a)
	if err != nil {
		log.Panic().Msg("middlewares not set up")
	}

	r.Use(middleware.Log(), gin.Recovery())
	r.GET("/check", check)
	r.POST("/api/register", h.Register)
	r.POST("/api/login", h.Login)
	r.POST("/api/companies", m.Auth(h.RegisterCompany))
	r.GET("/api/companies", m.Auth(h.GetCompanies))
	r.GET("/api/companies/:companyID", m.Auth(h.GetCompany))
	r.POST("/api/companies/:companyID/jobs", m.Auth(h.AddJobToCompany))
	r.GET("/api/companies/:companyID/jobs", m.Auth(h.ListCompanyJobs))
	r.GET("/api/jobs", m.Auth(h.ListJobs))
	r.GET("/api/jobs/:jobID", m.Auth(h.GetJobByID))
	return r
}

func check(c *gin.Context) {
	select {
	case <-c.Request.Context().Done():
		log.Info().Msg("User not there")
		return
	default:
		c.JSON(http.StatusOK, gin.H{"msg": http.StatusText(http.StatusOK)})
	}
}
