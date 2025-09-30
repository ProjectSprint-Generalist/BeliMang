package dto

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required,min=5,max=30"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}
