package usecases

import (
	"errors"
	"log"

	"findMyDoc/internal/entities"
	"findMyDoc/pkg/auth"
	"findMyDoc/pkg/email"
	"findMyDoc/users/repositories"
	"findMyDoc/users/requests"

	"golang.org/x/crypto/bcrypt"

	"math/rand"
	"time"
)

type UserUsecase interface {
	LoginUser(email, password string) (string, error)
	RegisterUser(req requests.RegisterRequest) error
}

type userUsecase struct {
	repo         repositories.UserRepository
	emailService email.EmailService
	appHost      string
}

func NewUserUsecase(repo repositories.UserRepository, emailService email.EmailService, appHost string) UserUsecase {
	return &userUsecase{
		repo:         repo,
		emailService: emailService,
		appHost:      appHost,
	}
}

func (uc *userUsecase) RegisterUser(req requests.RegisterRequest) error {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	verificationCode := generateVerificationCode()

	// Create user entity
	user := entities.User{
		Email:            req.Email,
		Password:         string(hashedPassword),
		Role:             req.Role,
		IsVerified:       false,
		VerificationCode: verificationCode,
	}

	// Save user to the database
	if err := uc.repo.CreateUser(&user); err != nil {
		return errors.New("failed to create user")
	}

	// Create doctor or patient profile
	if req.Role == "doctor" {
		doctor := entities.Doctor{
			ID:         int(user.ID),
			Name:       req.Name,
			Speciality: req.Speciality,
			Latitude:   req.Latitude,
			Longitude:  req.Longitude,
			Ratings:    0.0,
		}
		if err := uc.repo.CreateDoctor(&doctor); err != nil {
			return errors.New("failed to create doctor profile")
		}
	} else if req.Role == "patient" {
		patient := entities.Patient{
			ID:      int(user.ID),
			Name:    req.Name,
			Ratings: 0.0,
		}
		if err := uc.repo.CreatePatient(&patient); err != nil {
			return errors.New("failed to create patient profile")
		}
	} else {
		return errors.New("invalid role")
	}

	// Send verification email
	if err := uc.emailService.SendVerificationEmail(user.Email, verificationCode); err != nil {
		log.Printf("Failed to send email: %v", err) // <- Add this
		return errors.New("failed to send verification email")
	}
	return err
}

// LoginUser verifies the password and returns a JWT token
func (u *userUsecase) LoginUser(email, password string) (string, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	return auth.GenerateToken(int(user.ID), user.Role)
}

func generateVerificationCode() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
