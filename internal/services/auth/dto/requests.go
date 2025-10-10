package dto

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=5,max=30"`
	Password string `json:"password" binding:"required,min=5,max=30"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required,min=5,max=30"`
	Password string `json:"password" binding:"required,min=5,max=30"`
}
