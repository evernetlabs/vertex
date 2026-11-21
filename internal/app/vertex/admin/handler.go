package admin

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

	h.router.POST("/api/v1/admins/init", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		var req InitRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(&exception.ClientError{Message: err.Error()})
			return
		}

		admin, err := h.manager.Init(ctx, &req)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, admin)
	})

	h.router.POST("/api/v1/admins/token", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		var req TokenRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(&exception.ClientError{Message: err.Error()})
			return
		}

		token, err := h.manager.GetToken(ctx, &req)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, token)
	})

	h.router.GET("/api/v1/admins/current", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		aa, err := h.authenticator.ValidateAdminContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		admin, err := h.manager.Get(ctx, aa.Username)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, admin)
	})

	h.router.PUT("/api/v1/admins/current/password", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		aa, err := h.authenticator.ValidateAdminContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		var req PasswordChangeRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(&exception.ClientError{Message: err.Error()})
			return
		}

		err = h.manager.ChangePassword(ctx, aa.Username, &req)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	})

	h.router.POST("/api/v1/admins", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		aa, err := h.authenticator.ValidateAdminContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		var req AdditionRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(&exception.ClientError{Message: err.Error()})
			return
		}

		resp, err := h.manager.Add(ctx, &req, aa.Username)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, resp)
	})

	h.router.GET("/api/v1/admins", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateAdminContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		page, size := api.Page(c)

		admins, err := h.manager.List(ctx, page, size)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, admins)
	})

	h.router.GET("/api/v1/admins/:username", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateAdminContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		username := c.Param("username")

		admin, err := h.manager.Get(ctx, username)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, admin)
	})

	h.router.PUT("/api/v1/admins/:username/password", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateAdminContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		username := c.Param("username")

		resp, err := h.manager.ResetPassword(ctx, username)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusOK, resp)
	})

	h.router.DELETE("/api/v1/admins/:username", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateAdminContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		username := c.Param("username")

		err = h.manager.Delete(ctx, username)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	})
}
