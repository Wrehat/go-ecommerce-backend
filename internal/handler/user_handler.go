package handler

import (
	"ecommerce/internal/domain"
	"ecommerce/pkg/response"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RegisterReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type userHandler struct {
	usecase domain.UserUsecase
}

func NewUserHandler(uh domain.UserUsecase) *userHandler {
	return &userHandler{usecase: uh}
}

// Register godoc
// @Summary Register User Baru
// @Description Mendaftarkan akun pelanggan baru ke sistem
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RegisterReq true "Format JSON untuk register"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 409 {object} response.APIResponse
// @Router /users/register [post]
func (h *userHandler) Register(c *gin.Context) {
	// Todo : terima json dan validasi input
	var req RegisterReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "Format Input tidak valid", nil)
		return
	}

	// Todo : ubah jadi domain user
	user := domain.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
	}

	// Todo : panggil usecase
	userReg, err := h.usecase.Register(c.Request.Context(), user)

	// Todo : handle error
	if err != nil {
		if errors.Is(err, domain.ErrEmailDuplicate) {
			response.JSON(c, http.StatusConflict, "Email sudah digunakan", nil)
			return
		}
		response.JSON(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
		return
	}

	// Todo : return response
	response.JSON(c, http.StatusCreated, "Registrasi berhasil", userReg)
}

// Login godoc
// @Summary Login User
// @Description Melakukan autentikasi dan mendapatkan Token JWT
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body LoginReq true "Format JSON untuk login"
// @Success 200 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Failure 401 {object} response.APIResponse
// @Router /users/login [post]
func (h *userHandler) Login(c *gin.Context) {
	// Todo : terima JSON & validasi
	var req LoginReq

	if err := c.ShouldBindJSON(&req); err != nil {
		response.JSON(c, http.StatusBadRequest, "Format input tidak valid", nil)
		return
	}

	// Todo : panggil usecase dan func Login
	token, err := h.usecase.Login(c.Request.Context(), req.Email, req.Password)

	// Todo : handle error
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.JSON(c, http.StatusUnauthorized, "Email atau password anda salah", nil)
			return
		}
		response.JSON(c, http.StatusInternalServerError, "Terjadi kesalahan pada server", nil)
		return
	}

	// Todo : return response
	response.JSON(c, http.StatusOK, "Berhasil login", token)
}
