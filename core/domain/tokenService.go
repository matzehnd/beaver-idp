package domain

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenService struct {
	privateKey []byte
}

func NewTokenService(privateKey []byte) *TokenService {
	return &TokenService{
		privateKey: privateKey,
	}
}

func (s *TokenService) CreateToken(sub string, validityInHours *int, isAdmin bool) (string, error) {
	token, err := tokenFromUser(sub, isAdmin, validityInHours, s.privateKey)

	if err != nil {
		return "", fmt.Errorf("unable to get token: %T", err)
	}

	return token, nil
}

func tokenFromUser(sub string, isAdmin bool, validityInHours *int, privateKey []byte) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)

	duration := func(validity *int) time.Duration {
		if validity != nil {
			return time.Duration(time.Hour * time.Duration(*validity))
		}
		return time.Duration(time.Hour * 72)
	}(validityInHours)
	fmt.Println(duration)

	if err != nil {
		return "", err
	}
	claims := jwt.MapClaims{
		"isAdmin": isAdmin,
		"sub":     sub,
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *TokenService) ValidateToken(token string) (*jwt.Token, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return &key.PublicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("unable to verify token: %v", err)
	}

	return parsed, nil
}
