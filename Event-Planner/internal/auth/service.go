package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"event-planner/internal/user"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req user.RegisterRequest) (*user.AuthResponse, error) {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Insert user into database
	var u user.User
	query := `INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id, email, created_at`
	err = s.db.QueryRow(ctx, query, req.Email, string(hashedPassword)).Scan(&u.ID, &u.Email, &u.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.generateToken(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &user.AuthResponse{
		Token: token,
	}, nil
}

// Login authenticates a user and returns a token
func (s *Service) Login(ctx context.Context, req user.LoginRequest) (*user.AuthResponse, error) {
	var u user.User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	err := s.db.QueryRow(ctx, query, req.Email).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.generateToken(u.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Clear password hash before returning
	u.PasswordHash = ""

	return &user.AuthResponse{
		Token: token,
	}, nil
}

// generateToken creates a JWT token for the user
func (s *Service) generateToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(), // 7 days
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Default for development
	}

	return token.SignedString([]byte(jwtSecret))
}

// ValidateToken validates and parses a JWT token
func (s *Service) ValidateToken(tokenString string) (int, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Default for development
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid user ID in token")
		}
		return int(userID), nil
	}

	return 0, errors.New("invalid token")
}
