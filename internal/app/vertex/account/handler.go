package account

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

	h.router.POST("/api/v1/accounts/signup", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		var req SignUpRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			_ = c.Error(&exception.ClientError{Message: err.Error()})
			return
		}

		account, err := h.manager.SignUp(ctx, &req)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.JSON(http.StatusCreated, account)
	})

	h.router.POST("/api/v1/accounts/token", func(c *gin.Context) {
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
