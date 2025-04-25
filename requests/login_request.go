package requests

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LogoutRequest struct {
	AuthToken string `json:"auth_token" binding:"required"`
}
