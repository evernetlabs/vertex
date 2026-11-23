package profile

import (
	"context"
	"net/http"
	"time"

	"github.com/evernetlabs/vertex/internal/app/vertex/auth"
	"github.com/evernetlabs/vertex/internal/pkg/api"
	"github.com/evernetlabs/vertex/internal/pkg/exception"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	router        *gin.Engine
	authenticator *auth.Authenticator
	manager       *Manager
}

func NewHandler(router *gin.Engine, authenticator *auth.Authenticator, manager *Manager) *Handler {
	return &Handler{router: router, authenticator: authenticator, manager: manager}
}

func (h *Handler) Register() {

	h.router.POST("/api/v1/profiles", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		aa, err := h.authenticator.ValidateAccountContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		var req CreationRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(&exception.ClientError{Message: err.Error()})
			return
		}

		profile, err := h.manager.Create(ctx, &req, aa.Username)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, profile)
	})

	h.router.GET("/api/v1/profiles", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		aa, err := h.authenticator.ValidateAccountContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		page, size := api.Page(c)

		profiles, err := h.manager.List(ctx, aa.Username, page, size)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, profiles)
	})

	h.router.POST("/api/v1/profiles/:handle/token", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		aa, err := h.authenticator.ValidateAccountContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		handle := c.Param("handle")

		token, err := h.manager.GetToken(ctx, handle, aa.Username)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, token)
	})

	h.router.GET("/api/v1/profiles/current", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		ap, err := h.authenticator.ValidateProfileContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		profile, err := h.manager.Get(ctx, ap.Handle, ap.AccountUsername)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, profile)
	})

	h.router.PUT("/api/v1/profiles/current", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		ap, err := h.authenticator.ValidateProfileContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		var req UpdateRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(&exception.ClientError{Message: err.Error()})
			return
		}

		err = h.manager.Update(ctx, ap.Handle, &req, ap.AccountUsername)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	})

	h.router.DELETE("/api/v1/profiles/current", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		ap, err := h.authenticator.ValidateProfileContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		err = h.manager.Delete(ctx, ap.Handle, ap.AccountUsername)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	})
}
