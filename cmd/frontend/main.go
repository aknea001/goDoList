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

		switch strings.ToLower(command) {
		case "new":
			newTask := pkg.Task{}

			fmt.Print("Title: ")
			titleScanner := bufio.NewScanner(os.Stdin)
			titleScanner.Scan()
			newTask.Title = titleScanner.Text()

			fmt.Print("Description: ")
			descScanner := bufio.NewScanner(os.Stdin)
			descScanner.Scan()
			newTask.Description = descScanner.Text()

			err = api.NewTask(newTask)
			if err != nil {
				log.Fatal(err)
			}

			tasks = append(tasks, newTask)
		case "fn":
			var id int

			if len(inputSlice) == 2 {
				id, err = strconv.Atoi(inputSlice[1])
				if err != nil {
					fmt.Println("ID has to be of type int")
					time.Sleep(1 * time.Second)

					continue mainloop
				}
			} else {
				fmt.Print("ID: ")
				fmt.Scan(&id)
			}

			err = api.FinishTask(tasks[id-1])
			if err != nil {
				log.Fatal(err)
			}

			tasks = append(tasks[:id-1], tasks[id:]...)
		default:
			fmt.Printf("%s not found\n", command)
		}
	}
}
