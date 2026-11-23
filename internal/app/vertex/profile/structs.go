package profile

type CreationRequest struct {
	Handle      string `json:"handle" binding:"required"`
	DisplayName string `json:"display_name" binding:"required"`
	Type        string `json:"type" binding:"required"`
}

type UpdateRequest struct {
	DisplayName string `json:"display_name"`
	Type        string `json:"type"`
}

type TokenResponse struct {
	Token string `json:"token"`
}
