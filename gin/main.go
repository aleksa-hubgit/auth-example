package main

import (
	v1 "example-grpc/gen/proto/hello"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var jwtKey = []byte("my_secret_key")
var tokens []string

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := v1.NewGreeterClient(conn)

	r := gin.Default()
	r.GET("/rest/n/:name", func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")
		reqToken := strings.Split(bearerToken, " ")[1]
		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(reqToken, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "unauthorized",
				})
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			return
		}
		if !tkn.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "unauthorized",
			})
			return
		}
		name := c.Param("name")

		req := &v1.HelloRequest{Name: name}
		res, err := client.SayHello(c, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": fmt.Sprint(res.Mesasge),
		})
	})
	r.POST("/login", gin.BasicAuth(gin.Accounts{
		"admin": "secret",
	}), func(c *gin.Context) {
		token, err := generateJWT()
		if err != nil {
			fmt.Printf("failed to create jwt token: %v", err)
			return
		}
		// tokens := append(tokens, token)

		c.JSON(http.StatusOK, gin.H{
			"token": token,
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("could not run server: %v", err)
	}

}

func generateJWT() (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: "username",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	return token.SignedString(jwtKey)
}
