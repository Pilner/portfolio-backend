package token

import (
	"errors"
	"frv-backend/internal/adapters/config"
	"frv-backend/internal/core/domain"
	authdomain "frv-backend/internal/core/domain/auth"
	tokendomain "frv-backend/internal/core/domain/token"
	"frv-backend/internal/core/ports"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JwtTokenManager struct {
	config config.Values
}

func NewJwtTokenManager(cfg config.Values) ports.TokenManager {
	return JwtTokenManager{config: cfg}
}

func (j JwtTokenManager) GenerateToken(tokenType tokendomain.TokenType, payload authdomain.User) (string, error) {
	var key []byte
	var duration int
	if tokenType == tokendomain.TokenTypeJwt {
		key = []byte(j.config.JwtSecretKey)
		duration = j.config.JwtTokenExpiryMinutes
	}
	if tokenType == tokendomain.TokenTypeRefresh {
		key = []byte(j.config.RefreshTokenSecretKey)
		duration = j.config.RefreshTokenExpiryMinutes
	}

	claims := tokendomain.TokenClaim{
		User: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(duration) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "frv-backend",
			Subject:   payload.Id,
			Audience:  jwt.ClaimStrings{payload.Email},
			ID:        uuid.NewString(),
		},
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	jwtString, err := t.SignedString(key)
	if err != nil {
		return jwtString, err
	}

	return jwtString, nil
}

func (j JwtTokenManager) ValidateToken(tokenType tokendomain.TokenType, tokenString string) (*authdomain.User, error) {
	var key []byte
	if tokenType == tokendomain.TokenTypeJwt {
		key = []byte(j.config.JwtSecretKey)
	}
	if tokenType == tokendomain.TokenTypeRefresh {
		key = []byte(j.config.RefreshTokenSecretKey)
	}

	t, err := jwt.ParseWithClaims(tokenString, &tokendomain.TokenClaim{}, func(tok *jwt.Token) (any, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return key, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, domain.ErrExpiredToken
		}

		return nil, err
	}

	if claims, ok := t.Claims.(*tokendomain.TokenClaim); ok && t.Valid {
		return &claims.User, nil
	}

	return nil, domain.ErrInvalidToken
}
