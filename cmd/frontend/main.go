package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aknea001/goDoList/pkg"
	"github.com/aknea001/goDoList/pkg/frontend"
)

func getCreds() (string, string) {
	var username string
	var passwd string

	fmt.Print("Username: ")
	fmt.Scan(&username)

	fmt.Print("Password: ")
	fmt.Scan(&passwd)

	return username, passwd
}

func main() {
	api := frontend.NewAPIconn("http://localhost:8080")

	var x int

loginLoop:
	for {
		fmt.Print("What do you want to do?\n[ 1 ] Login \n[ 2 ] Register\n[ 1/2 ]: ")
		fmt.Scan(&x)

		switch x {
		case 1:
			username, passwd := getCreds()

			err := api.Login(username, passwd)
			if err != nil {
				var credE *pkg.CredentialError
				if errors.As(err, &credE) {
					fmt.Println("wrong username or password")
					continue loginLoop
				}

				log.Fatal(err)
			}

			break loginLoop
		case 2:
			username, passwd := getCreds()

			err := api.Register(username, passwd)
			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Successfully registered %s\n\n", username)
			time.Sleep(1 * time.Second)
		}
	}

	tasks, err := api.GetTasks()
mainloop:
	for {
		if err != nil {
			log.Fatal(err)
		}

		frontend.DrawTable(tasks)
		fmt.Print(">> ")

		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		input := scanner.Text()

		inputSlice := strings.Split(input, " ")
		command := inputSlice[0]

		switch command {
		case "fn":
			if len(inputSlice) == 2 {
				log.Print(inputSlice)
				id, err := strconv.Atoi(inputSlice[1])
				if err != nil {
					fmt.Println("ID has to be of type int")
					time.Sleep(1 * time.Second)

					continue mainloop
				}

				//send api request

				tasks = append(tasks[:id], tasks[id+1:]...)
			}

		}
	}
}
