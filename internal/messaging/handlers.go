package messaging

import (
	"context"
	"encoding/json"
	"log"
	"music-player-go/internal/auth"
	"music-player-go/pkg/models"

	"github.com/nats-io/nats.go"
)

func HandleRegistration(msg *nats.Msg) {
	var user models.User
	if err := json.Unmarshal(msg.Data, &user); err != nil {
		log.Println("Error unmarshaling registration data:", err)
		return
	}
	if err := auth.RegisterUser(context.Background(), &user); err != nil {
		log.Println("Error registering user:", err)
	}

}

func HandleLogin(msg *nats.Msg) {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.Unmarshal(msg.Data, &creds); err != nil {
		log.Println("Error unmarshaling login data:", err)
		return
	}

	token, err := auth.LoginUser(context.Background(), creds.Username, creds.Password)
	if err != nil {
		log.Println("Login failed:", err)
		return
	}

	log.Println("Login successful, token:", token)
}
