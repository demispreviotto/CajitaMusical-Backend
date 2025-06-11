package services

import (
	"context"
	"errors"
	"log" // Temporal para logging, en un proyecto grande se usaría un logger estructurado

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/dto/user"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm" // Para manejar errores específicos de GORM como record not found
)

// UserServicer define la interfaz para las operaciones del servicio de usuarios.
//
//go:generate mockgen -source=user_service.go -destination=mocks/mock_user_service.go
type UserServicer interface {
	RegisterUser(ctx context.Context, input user.RegisterUserInput) (*user.UserResponse, error)
	GetUserByID(ctx context.Context, userID string) (*user.UserInfo, error) // Asume que userID es string para compatibilidad inicial, luego cambiar a uuid.UUID
	// Agrega otros métodos de servicio de usuario aquí (ej. UpdateUser, DeleteUser)
}

// userService es la implementación concreta de UserServicer.
type userService struct {
	userDB db.UserDBer // Dependencia de la interfaz de la capa DB
}

// NewUserService crea una nueva instancia de UserService.
func NewUserService(userDB db.UserDBer) UserServicer {
	return &userService{userDB: userDB}
}

// RegisterUser maneja la lógica de negocio para registrar un nuevo usuario.
func (s *userService) RegisterUser(ctx context.Context, input user.RegisterUserInput) (*user.UserResponse, error) {
	// Aquí podrías añadir más validaciones de negocio antes de interactuar con la DB
	// Por ejemplo:
	// if !isValidUsername(input.Username) {
	//     return nil, errors.New("invalid username format")
	// }

	// 1. Hashear la contraseña
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Service: Failed to hash password for %s: %v", input.Username, err)
		return nil, errors.New("failed to process password")
	}

	// 2. Crear el modelo de usuario
	userModel := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Name:     input.Name,
	}

	// 3. Crear el usuario en la DB
	err = s.userDB.CreateUser(ctx, userModel, string(hashedPassword))
	if err != nil {
		// Aquí puedes manejar errores específicos de la DB, como usuario/email ya existente
		if errors.Is(err, gorm.ErrDuplicatedKey) { // Ejemplo, puede variar según el driver DB
			return nil, errors.New("username or email already registered")
		}
		log.Printf("Service: Failed to create user %s in DB: %v", input.Username, err)
		return nil, errors.New("failed to register user")
	}

	// 4. Mapear el modelo de DB a DTO de respuesta
	response := &user.UserResponse{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Name:     userModel.Name,
	}

	return response, nil
}

// GetUserByID maneja la lógica de negocio para obtener un usuario por ID.
func (s *userService) GetUserByID(ctx context.Context, userID string) (*user.UserInfo, error) {
	// Convertir el string ID a uuid.UUID si es necesario (asumimos que la DB espera uuid.UUID)
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	userModel, err := s.userDB.GetUserByID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		log.Printf("Service: Failed to get user by ID %s from DB: %v", userID, err)
		return nil, errors.New("failed to retrieve user")
	}

	// Mapear el modelo de DB a DTO de información de usuario
	userInfo := &user.UserInfo{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Name:     userModel.Name,
	}

	return userInfo, nil
}
