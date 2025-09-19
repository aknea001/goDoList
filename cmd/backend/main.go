package main

import (
	"errors"
	"log"

	"github.com/aknea001/goDoList/pkg"
	"github.com/aknea001/goDoList/pkg/backend"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/register", func(ctx *gin.Context) {
		var UserData pkg.User

		err := ctx.ShouldBindJSON(&UserData)
		if err != nil {
			log.Fatal(err)
		}

		err = backend.RegisterJson(UserData.Username, UserData.Passwd)
		if err != nil {
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

			log.Fatal(err)
		}

		ctx.JSON(200, gin.H{
			"msg":   "Success",
			"token": nil,
		})
	})

	router.Run()
}
