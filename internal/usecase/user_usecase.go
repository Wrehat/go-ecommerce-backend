package usecase

import (
	"context"
	"ecommerce/internal/domain"
	"ecommerce/pkg/jwtutil"
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
	// Todo : buat hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)

	if err != nil {
		return domain.User{}, err
	}

	user.PasswordHash = string(hashedPassword)

	// Todo : panggil repo dan buat data
	createdUser, err := u.repo.CreateUser(ctx, user)

	if err != nil {
		return domain.User{}, err
	}

	createdUser.PasswordHash = ""

	return createdUser, nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	// Todo : cari email di database
	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", domain.ErrInvalidCredentials
		}
		return "", err
	}

	// Todo : cek password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return "", domain.ErrInvalidCredentials
	}

	// Todo : generate token
	claims := jwtutil.MyCustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(u.jwtSecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
