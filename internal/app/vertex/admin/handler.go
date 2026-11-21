package admin

import (
	"context"
	"net/http"
	"time"

	"github.com/evernetlabs/vertex/internal/app/vertex/auth"
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
}
