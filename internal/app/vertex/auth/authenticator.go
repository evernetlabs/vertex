package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/evernetlabs/vertex/internal/app/vertex/config"
	"github.com/evernetlabs/vertex/internal/pkg/exception"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Authenticator struct {
	configManager *config.Manager
}

func NewAuthenticator(configManager *config.Manager) *Authenticator {
	return &Authenticator{configManager: configManager}
}

const (
	BearerToken      = "Bearer"
	TokenTypeAdmin   = "admin"
	TokenTypeAccount = "account"
)

func (a *Authenticator) GenerateAdminToken(ctx context.Context, username string) (string, error) {
	vertexEndpoint := a.configManager.GetVertexEndpoint(ctx)
	jwtSigningKey := a.configManager.GetJwtSigningKey(ctx)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  username,
		"iss":  vertexEndpoint,
		"aud":  vertexEndpoint,
		"type": TokenTypeAdmin,
		"iat":  int(time.Now().Unix()),
		"exp":  int(time.Now().Add(time.Hour * 24).Unix()),
	})

	tokenString, err := token.SignedString([]byte(jwtSigningKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *Authenticator) ValidateAdminContext(ctx context.Context, c *gin.Context) (*AuthenticatedAdmin, error) {
	tokenType, token, err := a.extractToken(c)

	if err != nil {
		return nil, err
	}

	switch tokenType {
	case BearerToken:
		return a.validateAdminToken(ctx, token)
	default:
		return nil, fmt.Errorf("invalid token type")
	}
}

func (a *Authenticator) validateAdminToken(ctx context.Context, tokenString string) (*AuthenticatedAdmin, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.configManager.GetJwtSigningKey(ctx)), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"]
		if !ok {
			return nil, &exception.AuthenticationError{
				Message: "invalid token type",
			}
		}

		tokenTypeString, ok := tokenType.(string)
		if !ok {
			return nil, &exception.AuthenticationError{
				Message: "invalid token type",
			}
		}

		if tokenTypeString != TokenTypeAdmin {
			return nil, &exception.AuthenticationError{
				Message: "invalid token type",
			}
		}

		username, ok := claims["sub"]

		if !ok {
			return nil, &exception.AuthenticationError{
				Message: "invalid access token",
			}
		}

		usernameString, ok := username.(string)

		if !ok {
			return nil, &exception.AuthenticationError{
				Message: "invalid access token",
			}
		}

		return &AuthenticatedAdmin{
			Username: usernameString,
		}, nil
	} else {
		return nil, &exception.AuthenticationError{
			Message: "invalid access token",
		}
	}
}

func (a *Authenticator) GenerateAccountToken(ctx context.Context, username string) (string, error) {
	vertexEndpoint := a.configManager.GetVertexEndpoint(ctx)
	jwtSigningKey := a.configManager.GetJwtSigningKey(ctx)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":  username,
		"iss":  vertexEndpoint,
		"aud":  vertexEndpoint,
		"type": TokenTypeAccount,
		"iat":  int(time.Now().Unix()),
		"exp":  int(time.Now().Add(time.Hour * 24).Unix()),
	})

	tokenString, err := token.SignedString([]byte(jwtSigningKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (a *Authenticator) ValidateAccountContext(ctx context.Context, c *gin.Context) (*AuthenticatedAccount, error) {
	tokenType, token, err := a.extractToken(c)

	if err != nil {
		return nil, err
	}

	switch tokenType {
	case BearerToken:
		return a.validateAccountToken(ctx, token)
	default:
		return nil, fmt.Errorf("invalid token type")
	}
}

func (a *Authenticator) validateAccountToken(ctx context.Context, tokenString string) (*AuthenticatedAccount, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.configManager.GetJwtSigningKey(ctx)), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		tokenType, ok := claims["type"]
		if !ok {
			return nil, &exception.AuthenticationError{
				Message: "invalid token type",
			}
		}

		tokenTypeString, ok := tokenType.(string)
		if !ok {
			return nil, &exception.AuthenticationError{
				Message: "invalid token type",
			}
		}

		if tokenTypeString != TokenTypeAccount {
			return nil, &exception.AuthenticationError{
				Message: "invalid token type",
			}
		}

		username, ok := claims["sub"]

		if !ok {
			return nil, &exception.AuthenticationError{
				Message: "invalid access token",
			}
		}

		usernameString, ok := username.(string)

		if !ok {
			return nil, &exception.AuthenticationError{
				Message: "invalid access token",
			}
		}

		return &AuthenticatedAccount{
			Username: usernameString,
		}, nil
	} else {
		return nil, &exception.AuthenticationError{
			Message: "invalid access token",
		}
	}
}

func (a *Authenticator) extractToken(c *gin.Context) (string, string, error) {
	authorizationHeader := c.GetHeader("Authorization")

	if len(authorizationHeader) == 0 {
		return "", "", &exception.AuthenticationError{
			Message: "authorization header is not set",
		}
	}

	components := strings.Split(authorizationHeader, " ")

	if len(components) != 2 {
		return "", "", &exception.AuthenticationError{
			Message: "invalid access token",
		}
	}

	return components[0], components[1], nil
}
