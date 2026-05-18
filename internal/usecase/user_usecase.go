package usecase

import (
	"context"
	"ecommerce/internal/domain"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	repo         domain.UserRepository
	jwtSecretKey []byte
}

func NewUserUsecase(repo domain.UserRepository, secretKey string) domain.UserUsecase {
	return &userUsecase{
		repo:         repo,
		jwtSecretKey: []byte(secretKey),
	}
}

func (u *userUsecase) Register(ctx context.Context, user domain.User) (domain.User, error) {
	// Todo : buat hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)

	// Todo : error handling
	if err != nil {
		return domain.User{}, err
	}

	// Todo : timpa password dengan hashed password
	user.PasswordHash = string(hashedPassword)

	// Todo : panggil layer repo & create
	createdUser, err := u.repo.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	// Todo : cegah hashedpassword bocor diluar layer usecase
	createdUser.PasswordHash = ""

	// Todo : return response
	return createdUser, nil
}

func (u *userUsecase) Login(ctx context.Context, email string, password string) (string, error) {

	// Todo : cek apakah email ada pada DB
	user, err := u.repo.GetUserByEmail(ctx, email)

	// Todo : handling error
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	// Todo : validasi kecocokan password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Todo : generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	// Todo : signing token
	tokenString, err := token.SignedString(u.jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
