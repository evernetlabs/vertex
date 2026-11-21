package vertex

import (
	"context"
	"net/http"
	"time"

	"github.com/evernetlabs/vertex/internal/app/vertex/auth"
	"github.com/evernetlabs/vertex/internal/app/vertex/config"
	"github.com/evernetlabs/vertex/internal/pkg/secret"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	router        *gin.Engine
	authenticator *auth.Authenticator
	manager       *config.Manager
}

func NewHandler(router *gin.Engine, authenticator *auth.Authenticator, manager *config.Manager) *Handler {
	return &Handler{router: router, authenticator: authenticator, manager: manager}
}

func (h *Handler) Register() {

	h.router.GET("/api/v1/configs/vertex", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		vertexEndpoint := h.manager.GetVertexEndpoint(ctx)

		c.JSON(http.StatusOK, &Vertex{Endpoint: vertexEndpoint})
	})

	h.router.PUT("/api/v1/configs/jwt-signing-key", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateAdminContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		jwtSigningKey, err := secret.GenerateSecret(128)

		if err != nil {
			_ = c.Error(err)
			return
		}

		err = h.manager.SetJwtSigningKey(ctx, jwtSigningKey)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	})

	h.router.PUT("/api/v1/configs/vertex/endpoint", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 5*time.Second)
		defer cancel()

		_, err := h.authenticator.ValidateAdminContext(ctx, c)

		if err != nil {
			_ = c.Error(err)
			return
		}

		var request EndpointUpdateRequest

		if err := c.ShouldBindJSON(&request); err != nil {
			_ = c.Error(err)
			return
		}

		err = h.manager.SetVertexEndpoint(ctx, request.Endpoint)

		if err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusNoContent)
	})
}
