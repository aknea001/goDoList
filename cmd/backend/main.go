package main

import (
	"fmt"
	"log"

	"github.com/aknea001/goDoList/pkg/backend"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.POST("/register", func(ctx *gin.Context) {
		var UserData backend.User

		err := ctx.ShouldBindBodyWithJSON(&UserData)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(UserData.Hash)

		//backend.AddJsonUser()
		ctx.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.Run()
}
