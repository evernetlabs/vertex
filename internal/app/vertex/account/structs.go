package account

type SignUpRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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
