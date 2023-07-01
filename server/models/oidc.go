package models

// OIDC OAuth2
type GiteeUserInfo struct {
	Id    int    `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
}
