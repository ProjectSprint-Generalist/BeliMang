package dto

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=5,max=30"`
	Password string `json:"password" binding:"required,min=5,max=30"`
	Email    string `json:"email" binding:"required,email"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required,min=5,max=30"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}

type UserAuthResponse struct {
	Token string `json:"token"`
}
