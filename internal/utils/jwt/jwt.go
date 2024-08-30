package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"saranasistemsolusindo.com/gusen-admin/internal/constants"
)

type Claims struct {
	LoginID           string `json:"login_id"`
	OrderRestrictions string `json:"order_restrictions"`
	jwt.StandardClaims
}

var secretKey []byte

func generateSecureKey() string {
	// Generate a random 32-byte key
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal(err)
	}

	// Encode the key to a base64 string
	return base64.StdEncoding.EncodeToString(key)
}

func init() {
	// Generate a new secret key if not found in environment variables
	secretKeyStr := generateSecureKey()
	log.Println("Generated new secret key:", secretKeyStr)

	secretKey = []byte(secretKeyStr)
}

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		// Split the header to get the token part
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		// Parse the token
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check if the token is expired
		if claims.ExpiresAt < time.Now().Unix() {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// Check if OrderRestrictions is "q"
		if claims.OrderRestrictions != constants.IS_ADMIN {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		// Token is valid, proceed with the request
		next.ServeHTTP(w, r)
	})
}

func GenerateJWT(loginID, orderRestrictions string) (string, error) {
	// Define the expiration time
	expirationTime := time.Now().Add(2 * time.Hour)

	// Create the claims
	claims := &Claims{
		LoginID:           loginID,
		OrderRestrictions: orderRestrictions,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
