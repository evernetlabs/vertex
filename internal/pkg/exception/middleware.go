package exception

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		defer func() {
			if r := recover(); r != nil {
				c.AbortWithStatusJSON(
					http.StatusInternalServerError,
					ErrorResponse{
						Success: false,
						Message: "internal server error",
					},
				)
			}
		}()

		c.Next()

		// No errors
		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var (
			serverErr        *ServerError
			clientErr        *ClientError
			notFoundErr      *NotFoundError
			authErr          *AuthenticationError
			authorizationErr *AuthorizationError
		)

		switch {
		case errors.As(err, &clientErr):
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Message: clientErr.Message,
			})

		case errors.As(err, &notFoundErr):
			c.AbortWithStatusJSON(http.StatusNotFound, ErrorResponse{
				Success: false,
				Message: notFoundErr.Message,
			})

		case errors.As(err, &authErr):
			c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{
				Success: false,
				Message: authErr.Message,
			})

		case errors.As(err, &authorizationErr):
			c.AbortWithStatusJSON(http.StatusForbidden, ErrorResponse{
				Success: false,
				Message: authorizationErr.Message,
			})

		case errors.As(err, &serverErr):
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Message: serverErr.Message,
			})

		default:
			log.Println("unknown error", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{
				Success: false,
				Message: "internal server error",
			})
		}
	}
}
