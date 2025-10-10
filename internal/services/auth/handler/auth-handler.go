package handler

import (
	"errors"
	"net/http"

	base "github.com/ProjectSprint-Generalist/BeliMang/internal/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/domain"
	auth "github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/dto"
	"github.com/ProjectSprint-Generalist/BeliMang/internal/services/auth/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	h.register(c, domain.RoleAdmin)
}

func (h *AuthHandler) LoginAdmin(c *gin.Context) {
	h.login(c, domain.RoleAdmin)
}

func (h *AuthHandler) RegisterUser(c *gin.Context) {
	h.register(c, domain.RoleUser)
}

func (h *AuthHandler) LoginUser(c *gin.Context) {
	h.login(c, domain.RoleUser)
}

func (h *AuthHandler) register(c *gin.Context, role domain.Role) {
	var payload auth.RegisterRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, base.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please make sure you have provided a valid username, email, and password",
			Code:    http.StatusBadRequest,
		})
		return
	}

	// At this point, anything is already validated. If an error occurs, it must be because of duplicate user.
	token, err := h.svc.Register(c.Request.Context(), payload.Username, payload.Email, payload.Password, role)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, service.ErrDuplicateUser) {
			status = http.StatusConflict
		}
		c.JSON(status, base.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    status,
		})
		return
	}

	c.JSON(http.StatusCreated, auth.Response{
		Token: token,
	})
}

func (h *AuthHandler) login(c *gin.Context, role domain.Role) {
	var payload auth.LoginRequest

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, base.ErrorResponse{
			Success: false,
			Error:   "Invalid input: please make sure you have provided a valid username and password",
			Code:    http.StatusBadRequest,
		})
		return
	}

	token, err := h.svc.Login(c.Request.Context(), payload.Username, payload.Password, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, base.ErrorResponse{
			Success: false,
			Error:   err.Error(),
			Code:    http.StatusInternalServerError,
		})
		return
	}

	c.JSON(http.StatusOK, auth.Response{
		Token: token,
	})
}
