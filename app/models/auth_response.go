package models

type AuthResponse struct {
	Token       string       `json:"token"`
	User        *User        `json:"user"`
	Role        *Role        `json:"role"`
	Permissions []Permission `json:"permissions"`
}
