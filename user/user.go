package user

import (
	"math/rand"
	"strconv"
	"time"

	constant "blog/constant"
	db "blog/database"
	ent "blog/ent"

	jwt "github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

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

func CheckPassword(email, password string) (*ent.User, error) {
	user, err := db.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	err = CheckPasswordHash(password, user.HashedPassword)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GenerateToken(user ent.User) (string, string, error) {
	accessToken, err := GenerateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := GenerateRefreshToken(user)
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
func GenerateAccessToken(user ent.User) (string, error) {
	claims := CustomClaims{
		Email: user.Email,
		Role:  string(user.Role),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)), // access token expires in 15 minutes
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(accessTokenSecret)
}

// GenerateRefreshToken generates a refresh token
func GenerateRefreshToken(user ent.User) (string, error) {
	claims := CustomClaims{
		Email: user.Email,
		Role:  string(user.Role),
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
func CheckAccessToken(tokenString string) error {
	claims, err := DecodeToken(tokenString)
	if err != nil {
		return err
	}
	deadline := time.Unix(claims.ExpiresAt.Time.Unix(), 0)
	if time.Now().After(deadline) {
		return err
	}
	return nil
}

// check RefreshToken and if it is expired, return error
func CheckRefreshToken(tokenString string) error {
	claims, err := DecodeToken(tokenString)
	if err != nil {
		return err
	}
	deadline := time.Unix(claims.ExpiresAt.Time.Unix(), 0)
	if time.Now().After(deadline) {
		return err
	}
	return nil
}

// By using the refresh token, generate a new access token and refresh token
func RefreshToken(refreshToken string) (string, string, error) {
	claims, err := DecodeToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	user, err := db.GetUserByEmail(claims.Email)
	if err != nil {
		return "", "", err
	}

	accessToken, err := GenerateAccessToken(*user)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	// delete old token
	err = db.DeleteTokenByRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// save tokens to db
	err = db.SaveToken(user.Email, accessToken, newRefreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
