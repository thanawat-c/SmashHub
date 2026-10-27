package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"smashhub/internal/models"
)

type AuthService struct {
	db *gorm.DB
}

func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{db: db}
}

type AuthResponse struct {
	User  *models.User `json:"user"`
	Token string       `json:"token"`
}

// Generate simple secure session token
func generateToken() string {
	b := make([]byte, 24)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Register creates a new customer user
func (s *AuthService) Register(displayName, phone, email, password string) (*models.User, string, error) {
	displayName = strings.TrimSpace(displayName)
	phone = strings.TrimSpace(phone)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)

	if displayName == "" {
		return nil, "", errors.New("full name / display name is required")
	}
	if phone == "" {
		return nil, "", errors.New("phone number is required")
	}
	if len(password) < 4 {
		return nil, "", errors.New("password must be at least 4 characters")
	}

	// Check phone collision
	var existing models.User
	if err := s.db.Where("phone = ?", phone).First(&existing).Error; err == nil {
		return nil, "", errors.New("this phone number is already registered")
	}

	if email != "" {
		if err := s.db.Where("email = ?", email).First(&existing).Error; err == nil {
			return nil, "", errors.New("this email address is already registered")
		}
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	avatarURL := fmt.Sprintf("https://api.dicebear.com/7.x/avataaars/svg?seed=%s", phone)

	newUser := models.User{
		ID:           uuid.New(),
		DisplayName:  displayName,
		Phone:        phone,
		Email:        email,
		PictureURL:   avatarURL,
		PasswordHash: string(hashed),
		Role:         models.RoleCustomer,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(&newUser).Error; err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token := generateToken()
	return &newUser, token, nil
}

// Login authenticates with phone/email/lineID and password
func (s *AuthService) Login(identifier, password string) (*models.User, string, error) {
	identifier = strings.TrimSpace(identifier)
	password = strings.TrimSpace(password)

	if identifier == "" {
		return nil, "", errors.New("please provide a phone number, email, or username")
	}

	var user models.User
	err := s.db.Where("phone = ? OR email = ? OR line_user_id = ? OR display_name = ?",
		identifier, identifier, identifier, identifier).First(&user).Error

	if err != nil {
		return nil, "", errors.New("account not found. Please check your credentials or register.")
	}

	// If password hash exists, verify it
	if user.PasswordHash != "" {
		// Allow default demo password '123456' as fallback for recruiter convenience
		if password != "123456" {
			if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
				return nil, "", errors.New("incorrect password. Please try again.")
			}
		}
	}

	token := generateToken()
	return &user, token, nil
}

// DemoLogin provides fast 1-click login for test roles
func (s *AuthService) DemoLogin(role string) (*models.User, string, error) {
	var user models.User
	var err error

	if strings.ToUpper(role) == "ADMIN" {
		err = s.db.Where("role = ?", models.RoleAdmin).First(&user).Error
	} else {
		err = s.db.Where("role = ?", models.RoleCustomer).First(&user).Error
	}

	if err != nil {
		return nil, "", errors.New("test user for demo mode not found")
	}

	token := generateToken()
	return &user, token, nil
}

// GetMe retrieves user profile by UUID
func (s *AuthService) GetMe(userID string) (*models.User, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, errors.New("invalid user id format")
	}

	var user models.User
	if err := s.db.First(&user, "id = ?", uid).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
