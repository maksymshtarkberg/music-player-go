package cmd

import (
	"encoding/json"
	"log"

	"github.com/maksymshtarkberg/music-player-go/internal/auth"
	"github.com/maksymshtarkberg/music-player-go/internal/database"
	"github.com/maksymshtarkberg/music-player-go/internal/messaging"
	"github.com/maksymshtarkberg/music-player-go/pkg/models"

	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/nats-io/nats.go"
)

var nc *nats.Conn

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB_NAME")
	database.InitMongo(mongoURI, dbName)

	redisAddr := os.Getenv("REDIS_ADDR")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	auth.InitRedis(redisAddr, redisPassword, 0)

	natsURL := os.Getenv("NATS_URL")
	messaging.InitNATS(natsURL)
	nc = messaging.GetNATSConn()

	messaging.StartSubscribers()

	router := gin.Default()

	router.POST("/api/v1/user/reg", RegisterUserNats)
	router.POST("/api/v1/user/auth", AuthUserNats)

	router.Run(":8080")

	log.Println("Server is running...")
	select {}
}

func RegisterUserNats(c *gin.Context) {
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	jsonDataBytes, err := json.Marshal(jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := nc.Request("user.register", jsonDataBytes, nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response.Data, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}

func AuthUserNats(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userBytes, err := json.Marshal(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response, err := nc.Request("user.login", userBytes, nats.DefaultTimeout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(response.Data, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
