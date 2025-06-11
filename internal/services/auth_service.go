package services

import (
	"context"
	"errors"
	"log"
	"net"
	"os"
	"time"

	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/db"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/dto/auth"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/dto/user"
	"github.com/demispreviotto/cajitamusical/cajitamusical-backend/internal/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthServicer define la interfaz para las operaciones del servicio de autenticación.
//
//go:generate mockgen -source=auth_service.go -destination=mocks/mock_auth_service.go
type AuthServicer interface {
	Login(ctx context.Context, username, password, userAgent, ipAddress string) (*auth.LoginResponse, *models.Session, error)
	Logout(ctx context.Context, sessionID uuid.UUID) error
	CleanupExpiredSessions(ctx context.Context) error
}

// authService es la implementación concreta de AuthServicer.
type authService struct {
	userDB    db.UserDBer    // Dependencia de la interfaz de la capa DB de usuarios
	sessionDB db.SessionDBer // Dependencia de la interfaz de la capa DB de sesiones
}

// NewAuthService crea una nueva instancia de AuthService.
func NewAuthService(userDB db.UserDBer, sessionDB db.SessionDBer) AuthServicer {
	return &authService{userDB: userDB, sessionDB: sessionDB}
}

// Login maneja la lógica de negocio para el inicio de sesión.
func (s *authService) Login(ctx context.Context, username, password, userAgent, ipAddress string) (*auth.LoginResponse, *models.Session, error) {
	// 1. Obtener usuario y contraseña hasheada de la DB
	userModel, hashedPassword, err := s.userDB.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("invalid credentials")
		}
		log.Printf("Service: Failed to get user by username %s from DB: %v", username, err)
		return nil, nil, errors.New("login failed due to server error")
	}

	// 2. Comparar contraseñas
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return nil, nil, errors.New("invalid credentials")
	}

	// 3. Crear una nueva sesión
	sessionDurationStr := os.Getenv("SESSION_DURATION_HOURS")
	sessionDurationHours := 24 // Default
	if dur, errParse := time.ParseDuration(sessionDurationStr + "h"); errParse == nil {
		sessionDurationHours = int(dur.Hours())
	}
	expiresAt := time.Now().Add(time.Duration(sessionDurationHours) * time.Hour)

	// Convertir string de IP a net.IP
	clientIP := net.ParseIP(ipAddress)
	if clientIP == nil {
		log.Printf("Service: Warning: Could not parse client IP '%s'. Using nil/default.", ipAddress)
	}

	session := &models.Session{
		SessionID: uuid.New(), // Generar un nuevo UUID para la sesión
		UserID:    userModel.ID,
		CreatedAt: time.Now(),
		ExpiresAt: expiresAt,
		UserAgent: userAgent,
		IPAddress: clientIP,
	}

	if err := s.sessionDB.CreateSession(ctx, session); err != nil {
		log.Printf("Service: Failed to create session for user %s: %v", userModel.ID.String(), err)
		return nil, nil, errors.New("failed to create session")
	}

	// 4. Mapear el modelo de usuario a DTO de respuesta
	responseUser := user.UserInfo{
		ID:       userModel.ID,
		Username: userModel.Username,
		Email:    userModel.Email,
		Name:     userModel.Name,
	}

	loginResponse := &auth.LoginResponse{
		Message: "Login successful",
		User:    responseUser,
	}

	return loginResponse, session, nil
}

// Logout maneja la lógica de negocio para cerrar sesión.
func (s *authService) Logout(ctx context.Context, sessionID uuid.UUID) error {
	err := s.sessionDB.DeleteSession(ctx, sessionID)
	if err != nil {
		log.Printf("Service: Failed to delete session %s: %v", sessionID.String(), err)
		return errors.New("failed to logout fully")
	}
	return nil
}

// CleanupExpiredSessions maneja la lógica de negocio para limpiar sesiones expiradas.
func (s *authService) CleanupExpiredSessions(ctx context.Context) error {
	err := s.sessionDB.DeleteExpiredSessions(ctx)
	if err != nil {
		log.Printf("Service: Error during expired session cleanup: %v", err)
		return errors.New("failed to cleanup expired sessions")
	}
	return nil
}
