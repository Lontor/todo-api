package service

import (
	"context"
	"net/http"
	"time"

	"github.com/Lontor/todo-api/internal/dto"
	"github.com/Lontor/todo-api/internal/model"
	"github.com/Lontor/todo-api/internal/repository"
	"github.com/Lontor/todo-api/pkg/custom_errors"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	r repository.UserRepository
	v *validator.Validate
}

func NewUserService(r repository.UserRepository) UserService {
	return &userService{
		r: r,
		v: validator.New(),
	}
}

func (s *userService) CreateUser(ctx context.Context, data dto.RegisterRequest) error {
	if err := s.v.Struct(data); err != nil {
		return custom_errors.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(data.Password), 14)
	if err != nil {
		return custom_errors.NewHTTPError(http.StatusInternalServerError, "hash generation error")
	}

	if data.Role == "" {
		data.Role = model.UserTypeRegular
	}

	if role, ok := ctx.Value("role").(model.UserType); ok {
		if role != model.UserTypeAdmin {
			return custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
		return s.r.Create(ctx, model.User{
			ID:           uuid.New(),
			Email:        data.Email,
			PasswordHash: string(passwordHash),
			AccountType:  data.Role,
		})
	}

	if data.Role != model.UserTypeRegular {
		return custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
	}

	return s.r.Create(ctx, model.User{
		ID:           uuid.New(),
		Email:        data.Email,
		PasswordHash: string(passwordHash),
		AccountType:  data.Role,
	})
}

func (s *userService) GetUsers(ctx context.Context) ([]model.User, error) {
	role, ok := ctx.Value("role").(model.UserType)
	if !ok || role != model.UserTypeAdmin {
		return nil, custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
	}
	return s.r.Get(ctx)
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (model.User, error) {
	userID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	if role != model.UserTypeAdmin {
		if id != userID {
			return model.User{}, custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	return s.r.GetByID(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, data dto.UpdateUserRequest) error {
	if err := s.v.Struct(data); err != nil {
		return custom_errors.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	userID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	if role != model.UserTypeAdmin {
		if userID != data.UserID || data.Role == model.UserTypeAdmin {
			return custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	if data.Email == "" && data.Password == "" && data.Role == "" {
		return custom_errors.NewHTTPError(http.StatusBadRequest, "no fields to update")
	}

	var passwordHash []byte
	if data.Password != "" {
		var err error
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(data.Password), 14)
		if err != nil {
			return custom_errors.NewHTTPError(http.StatusInternalServerError, "hash generation error")
		}
	}

	return s.r.Update(ctx, model.User{
		ID:           data.UserID,
		Email:        data.Email,
		PasswordHash: string(passwordHash),
		AccountType:  data.Role,
	})
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	userID := ctx.Value("userID").(uuid.UUID)
	role := ctx.Value("role").(model.UserType)

	if role != model.UserTypeAdmin {
		if id != userID {
			return custom_errors.NewHTTPError(http.StatusForbidden, "permission denied")
		}
	}

	return s.r.Delete(ctx, id)
}

func (s *userService) AuthenticateUser(ctx context.Context, email, password string) (dto.AuthResponse, error) {
	user, err := s.r.GetByEmail(ctx, email)
	if err != nil {
		return dto.AuthResponse{}, custom_errors.NewHTTPError(http.StatusNotFound, "user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return dto.AuthResponse{}, custom_errors.NewHTTPError(http.StatusUnauthorized, "invalid credentials")
	}

	expirationTime := time.Now().Add(time.Hour * 3)

	claims := jwt.MapClaims{
		"userID": user.ID.String(),
		"role":   user.AccountType,
		"exp":    expirationTime.Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte("Y80zXN/dnoc14mdIpchh4ZXOdDAZfWulff5jEcCWHEc="))
	if err != nil {
		return dto.AuthResponse{}, custom_errors.NewHTTPError(http.StatusInternalServerError, "token generation error")
	}

	return dto.AuthResponse{
		Token:     tokenString,
		UserID:    user.ID,
		ExpiresAt: expirationTime,
	}, nil
}
