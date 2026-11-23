package auth

type AuthenticatedAdmin struct {
	Username string `json:"username"`
}

type AuthenticatedAccount struct {
	Username string `json:"username"`
}
