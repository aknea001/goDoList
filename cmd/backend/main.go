package main

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/aknea001/goDoList/pkg"
	"github.com/aknea001/goDoList/pkg/backend"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	router := gin.Default()
	router.POST("/register", func(ctx *gin.Context) {
		var UserData pkg.User

		err := ctx.ShouldBindJSON(&UserData)
		if err != nil {
			ctx.JSON(500, gin.H{
				"msg": "Unknown error",
			})

			log.Fatal(err)
		}

		err = backend.RegisterJson(UserData.Username, UserData.Passwd)
		if err != nil {
			ctx.JSON(500, gin.H{
				"msg": "Unknown error",
			})

			log.Fatal(err)
		}

		ctx.JSON(201, gin.H{
			"msg": "Success",
			"id":  nil,
		})
	})

	router.POST("/login", func(ctx *gin.Context) {
		var UserData pkg.User

		err := ctx.ShouldBindJSON(&UserData)
		if err != nil {
			ctx.JSON(500, gin.H{
				"msg": "Unknown error",
			})

			log.Fatal(err)
		}

		err = backend.LoginJson(UserData.Username, UserData.Passwd)
		if err != nil {
			var credE *pkg.CredentialError
			if errors.As(err, &credE) {
				ctx.JSON(401, gin.H{
					"msg": "Wrong username or password",
				})

				return
			}

			ctx.JSON(500, gin.H{
				"msg": "Unknown error",
			})

			log.Fatal(err)
		}

		jwtKey := []byte(os.Getenv("jwtKey"))

		// replace sub with IDs when proper DB implemented
		claims := &jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			Subject:   UserData.Username,
		}

		baseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		signedToken, err := baseToken.SignedString(jwtKey)
		if err != nil {
			ctx.JSON(500, gin.H{
				"msg": "Unknown error",
			})

			log.Fatal(err)
		}

		ctx.JSON(200, gin.H{
			"msg":   "Success",
			"token": signedToken,
		})
	})

	router.GET("/validateToken", func(ctx *gin.Context) {
		headers := ctx.Request.Header

		auth := headers.Get("Authorization")
		if auth == "" {
			ctx.JSON(401, gin.H{
				"msg": "Missing authorization of type Bearer",
			})
			return
		}

		authSplit := strings.Split(auth, " ")
		if len(authSplit) <= 1 || authSplit[0] != "Bearer" {
			ctx.JSON(401, gin.H{
				"msg": "Missing authorization of type Bearer",
			})

			return
		}

		signedToken := authSplit[1]
		validMethods := make([]string, 1)
		validMethods[0] = "HS256"

		token, err := jwt.Parse(signedToken, func(token *jwt.Token) (any, error) {
			return []byte(os.Getenv("jwtKey")), nil
		}, jwt.WithValidMethods(validMethods))

		if err != nil {
			ctx.JSON(500, gin.H{
				"msg": "Unknown error",
			})

			log.Fatal(err)
		}

		if !token.Valid {
			ctx.JSON(401, gin.H{
				"msg": "Invalid authorization",
			})

			return
		}

		ctx.JSON(200, gin.H{
			"msg": "Token is valid",
		})
	})

	router.Run()
}
