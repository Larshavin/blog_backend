package user

import (
	"math/rand"
	"strconv"
	"time"

	constant "blog/constant"
	es "blog/elastic-search"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	min_value, _ := strconv.Atoi(constant.HASH_MIN)
	max_value, _ := strconv.Atoi(constant.HASH_MAX)
	min := min_value
	max := max_value

	randomNumber := rand.Intn(max-min) + min
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), randomNumber)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func CheckPassword(email, password string, client es.Client) (es.User, error) {
	user, err := client.GetUserByEmail(email)
	if err != nil {
		return es.User{}, err
	}

	err = CheckPasswordHash(password, user.Password)
	if err != nil {
		return es.User{}, err
	}

	return user, nil
}

func GenerateToken(user es.User, client es.Client) (string, string, error) {
	accessToken, err := GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(user)
	if err != nil {
		return "", "", err
	}

	// save tokens to elastic search
	err = client.SaveToken(user.Email, accessToken, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Secret keys for signing tokens (in a real application, these should be stored securely)
var accessTokenSecret = []byte("access_secret_key")
var refreshTokenSecret = []byte("refresh_secret_key")

// Custom claims for Access and Refresh tokens
type CustomClaims struct {
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateAccessToken generates an access token
func GenerateAccessToken(user es.User) (string, error) {
	claims := CustomClaims{
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // access token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessTokenSecret)
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(user es.User) (string, error) {
	claims := CustomClaims{
		Email: user.Email,
		Role:  user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // refresh token expires in 1 days
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(refreshTokenSecret)
}

// decode token and get info
func DecodeToken(tokenString string) (CustomClaims, error) {
	claims := CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return accessTokenSecret, nil
	})
	if err != nil {
		return CustomClaims{}, err
	}
	if !token.Valid {
		return CustomClaims{}, err
	}
	return claims, nil
}

// check AccessToken and if it is expired, return error

// check RefreshToken and if it is expired, return error

// By using the refresh token, generate a new access token and refresh token
