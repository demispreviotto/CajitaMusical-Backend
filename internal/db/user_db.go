package db

import (
	"context"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
	"github.com/google/uuid" // Importar uuid
	"gorm.io/gorm"
)

// UserDBer define la interfaz para las operaciones de la base de datos de usuarios.
type UserDBer interface {
	CreateUser(ctx context.Context, user *models.User, hashedPassword string) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, string, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) // Usar uuid.UUID
	// Agrega otros métodos de DB de usuario aquí
}

// userDB es la implementación concreta de UserDBer.
type userDB struct{}

// NewUserDB crea una nueva instancia de UserDB.
func NewUserDB() UserDBer {
	return &userDB{}
}

// Implementación de CreateUser
func (udb *userDB) CreateUser(ctx context.Context, user *models.User, hashedPassword string) error {
	// Asegúrate de que user.ID se genere aquí si GORM no lo hace automáticamente
	// Por ejemplo: if user.ID == uuid.Nil { user.ID = uuid.New() }
	auth := models.Authentication{
		UserID:       user.ID,
		PasswordHash: hashedPassword,
	}

	// Inicia una transacción
	return DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		// Asegúrate de que el ID del usuario esté disponible para la autenticación
		auth.UserID = user.ID
		return tx.Create(&auth).Error
	})
}

// Implementación de GetUserByUsername
func (udb *userDB) GetUserByUsername(ctx context.Context, username string) (*models.User, string, error) {
	var user models.User
	var auth models.Authentication

	err := DB.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, "", err
	}

	err = DB.WithContext(ctx).Where("user_id = ?", user.ID).First(&auth).Error
	if err != nil {
		return nil, "", err // Si no encuentra auth, es un error interno o dato inconsistente
	}

	return &user, auth.PasswordHash, nil
}

// Implementación de GetUserByID
func (udb *userDB) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var user models.User
	err := DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	return &user, err
}
