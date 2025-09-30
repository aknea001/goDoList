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

func unknownError(err error, ctx *gin.Context) {
	ctx.JSON(500, gin.H{
		"msg": "unknown server error",
	})

	log.Print(err)
}

func validateToken(ctx *gin.Context) (string, error) {
	headers := ctx.Request.Header

	auth := headers.Get("Authorization")
	if auth == "" {
		ctx.JSON(401, gin.H{
			"msg": "Missing authorization of type Bearer",
		})

		return "", jwt.ErrTokenMalformed
	}

	authSplit := strings.Split(auth, " ")
	if len(authSplit) <= 1 || authSplit[0] != "Bearer" {
		ctx.JSON(401, gin.H{
			"msg": "Missing authorization of type Bearer",
		})

		return "", jwt.ErrTokenMalformed
	}

	signedToken := authSplit[1]
	validMethods := make([]string, 1)
	validMethods[0] = "HS256"

	token, err := jwt.Parse(signedToken, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("jwtKey")), nil
	}, jwt.WithValidMethods(validMethods))

	switch {
	case token.Valid:
		break

	// being a lil lazy, might make more in-dept res later
	case errors.Is(err, jwt.ErrTokenMalformed) ||
		errors.Is(err, jwt.ErrTokenSignatureInvalid) ||
		errors.Is(err, jwt.ErrTokenExpired):

		ctx.JSON(401, gin.H{
			"msg": "Invalid token",
		})

		return "", jwt.ErrTokenMalformed
	default:
		unknownError(err, ctx)

		return "", &pkg.UnknownServerError{}
	}

	sub, err := token.Claims.GetSubject()
	if err != nil {
		unknownError(err, ctx)

		return "", &pkg.UnknownServerError{}
	}

	return sub, nil
}

func main() {
	godotenv.Load()

	logFile, err := os.OpenFile(
		"goDoListBackend.log",
		os.O_WRONLY|os.O_APPEND|os.O_CREATE,
		0664,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	router := gin.Default()
	router.POST("/register", func(ctx *gin.Context) {
		var UserData pkg.User

		err := ctx.ShouldBindJSON(&UserData)
		if err != nil {
			unknownError(err, ctx)

			return
		}

		err = backend.RegisterJson(UserData.Username, UserData.Passwd)
		if err != nil {
			unknownError(err, ctx)

			return
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
			unknownError(err, ctx)

			return
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

			unknownError(err, ctx)

			return
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
			unknownError(err, ctx)

			return
		}

		ctx.JSON(200, gin.H{
			"msg":   "Success",
			"token": signedToken,
		})
	})

	router.GET("/tasks/get", func(ctx *gin.Context) {
		sub, err := validateToken(ctx)
		if err != nil {
			return
		}

		tasks, err := backend.GetTaskJson(sub)
		if err != nil {
			unknownError(err, ctx)
			return
		}

		ctx.JSON(200, gin.H{
			"msg":   "Success",
			"tasks": tasks,
		})
	})

	router.POST("/tasks/new", func(ctx *gin.Context) {
		sub, err := validateToken(ctx)
		if err != nil {
			return
		}

		var newTask pkg.Task

		err = ctx.ShouldBindJSON(&newTask)
		if err != nil {
			unknownError(err, ctx)
			return
		}

		newTask.Owner = sub

		err = backend.NewTaskJson(newTask)
		if err != nil {
			unknownError(err, ctx)
			return
		}

		ctx.JSON(200, gin.H{
			"msg": "Success",
		})
	})

	router.Run()
}
