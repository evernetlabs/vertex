package admin

import "github.com/evernetlabs/vertex/internal/app/vertex/vertex"

type InitRequest struct {
	Username string         `json:"username" binding:"required"`
	Password string         `json:"password" binding:"required"`
	Vertex   *vertex.Vertex `json:"vertex" binding:"required"`
}

type TokenRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type PasswordChangeRequest struct {
	Password string `json:"password" binding:"required"`
}

type AdditionRequest struct {
	Username string `json:"username" binding:"required"`
}

type AdditionResponse struct {
	Password string `json:"password"`
	Admin    *Admin `json:"admin"`
}

type PasswordResetResponse struct {
	Password string `json:"password"`
}
