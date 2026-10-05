package service

import (
	"authService/internal/domain"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type SecurityService struct {
	jwtSecret string
}

func NewSecurityService(jwtSecret string) *SecurityService {
	return &SecurityService{jwtSecret: jwtSecret}
}

func (service *SecurityService) HashPassword(password string) (string, error) {
	shaSum := sha256.Sum256([]byte(password))
	shaHex := hex.EncodeToString(shaSum[:])
	bytes, err := bcrypt.GenerateFromPassword([]byte(shaHex), bcrypt.DefaultCost)
	return string(bytes), err
}

func (service *SecurityService) GenerateTokens(userId string) (domain.Tokens, error) {
	expiresIn := time.Now().Add(15 * time.Minute)
	accessToken, err := service.generateToken(userId, expiresIn.Unix())
	if err != nil {
		return domain.Tokens{}, err
	}
	refreshExpiresIn := time.Now().Add(14 * 24 * time.Hour)
	refreshToken, err := service.generateToken(userId, refreshExpiresIn.Unix())
	if err != nil {
		return domain.Tokens{}, err
	}
	return domain.Tokens{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: expiresIn}, nil
}

func (service *SecurityService) generateToken(userId string, expiry int64) (string, error) {
	claims := jwt.MapClaims{
		"sub": userId,
		"exp": expiry,
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(service.jwtSecret))
}
