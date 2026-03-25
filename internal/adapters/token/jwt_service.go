package token

import (
	"errors"
	"portfolio-backend/internal/adapters/config"
	"portfolio-backend/internal/core/domain"
	authdomain "portfolio-backend/internal/core/domain/auth"
	tokendomain "portfolio-backend/internal/core/domain/token"
	"portfolio-backend/internal/core/ports"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
	config config.Values
}

func NewJwtService(cfg config.Values) ports.TokenService {
	return jwtService{config: cfg}
}

func (s jwtService) GenerateAccessToken(payload authdomain.User) (string, error) {
	key := []byte(s.config.JwtSecretKey)

	claims := tokendomain.TokenClaim{
		User: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.JwtTokenExpiryMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "portfolio-backend",
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

func (s jwtService) GenerateRefreshToken(payload authdomain.User) (string, error) {
	key := []byte(s.config.RefreshTokenSecretKey)

	claims := tokendomain.TokenClaim{
		User: payload,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(s.config.RefreshTokenExpiryMinutes) * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "portfolio-backend",
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

func (s jwtService) ValidateAccessToken(tokenString string) (*authdomain.User, error) {
	t, err := jwt.ParseWithClaims(tokenString, &tokendomain.TokenClaim{}, func(tok *jwt.Token) (any, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		key := []byte(s.config.JwtSecretKey)
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

func (s jwtService) ValidateRefreshToken(tokenString string) (*authdomain.User, error) {
	t, err := jwt.ParseWithClaims(tokenString, &tokendomain.TokenClaim{}, func(tok *jwt.Token) (any, error) {
		if _, ok := tok.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		key := []byte(s.config.RefreshTokenSecretKey)
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
