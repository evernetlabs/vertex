package auth

type AuthenticatedAdmin struct {
	Username string `json:"username"`
}

type AuthenticatedAccount struct {
	Username string `json:"username"`
}

type AuthenticatedProfile struct {
	Handle          string `json:"handle"`
	AccountUsername string `json:"account_username"`
}
