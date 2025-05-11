package main

import (
	"channel-service/internal/channel"
	"channel-service/internal/middleware"
	"channel-service/internal/subscription"
	"channel-service/pkg/database"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	// Define flags with default values (development-friendly)
	dbHost := flag.String("db-host", "localhost", "PostgreSQL host")
	dbPort := flag.String("db-port", "5432", "PostgreSQL port")
	dbUser := flag.String("db-user", "postgres", "PostgreSQL user")
	dbPassword := flag.String("db-password", "mysecretpassword", "PostgreSQL password")
	dbName := flag.String("db-name", "youtube", "PostgreSQL database name")
	jwtSecret := flag.String("jwt-secret", "dev_secret_123", "JWT signing secret")

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": "test_user_123",
	})
	tokenString, _ := token.SignedString([]byte("your_super_secret_key_here"))
	fmt.Println("Test token:", tokenString)
	// Initialize database
	db := database.NewPostgresDB(*dbHost, *dbPort, *dbUser, *dbPassword, *dbName)
	defer db.Close()

	// Setup handlers
	channelHandler := channel.NewHandler(db)
	subHandler := subscription.NewHandler(db)

	// Routes
	http.Handle("/channels/", middleware.AuthMiddleware(*jwtSecret, channelHandler))
	http.Handle("/subscriptions/", middleware.AuthMiddleware(*jwtSecret, subHandler))

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

}
