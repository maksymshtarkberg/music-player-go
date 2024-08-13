package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"music-player-go/internal/database"
	"music-player-go/pkg/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func CheckPassword(hash, password string) bool {
	return hash == HashPassword(password)
}

func RegisterUser(ctx context.Context, user *models.User) error {
	collection := database.GetCollection("users")

	user.Password = HashPassword(user.Password)

	_, err := collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errors.New("user already exists")
		}
		return err
	}
	return nil
}

func LoginUser(ctx context.Context, username, password string) (string, error) {

	collection := database.GetCollection("users")

	var user models.User
	err := collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !CheckPassword(user.Password, password) {
		return "", errors.New("invalid credentials")
	}

	existingToken, err := GetToken(ctx, user.ID)
	if err == nil && existingToken != "" {
		return "", errors.New("user already logged in")
	}

	expirationTime := time.Now().Add(time.Hour * 24 * 7)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	Key := os.Getenv("JWT_SECRET_KEY")
	var jwtKey = []byte(Key)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	if err := SetToken(ctx, tokenString, user.ID); err != nil {
		return "", err
	}

	return tokenString, nil
}
