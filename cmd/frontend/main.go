package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aknea001/goDoList/pkg"
	"github.com/aknea001/goDoList/pkg/frontend"
	"golang.org/x/term"
)

func getCreds() (string, string) {
	var username string

	fmt.Print("Username: ")
	fmt.Scan(&username)

	fmt.Print("Password: ")
	passwd, err := term.ReadPassword(1)
	if err != nil {
		log.Fatal(err)
	}

	return username, string(passwd)
}

func main() {
	api := frontend.NewAPIconn("http://localhost:8080")

	var x int
fullLoop:
	for {
	loginLoop:
		for {
			fmt.Print("What do you want to do?\n[ 1 ] Login \n[ 2 ] Register\n[ 1/2 ]: ")
			_, err := fmt.Scan(&x)
			if err != nil {
				fmt.Print("Please input either 1 or 2\n\n")
				fmt.Scanln()
				continue loginLoop
			}

			switch x {
			case 1:
				username, passwd := getCreds()

				err := api.Login(username, passwd)
				if err != nil {
					var credE *pkg.CredentialError
					if errors.As(err, &credE) {
						fmt.Print("wrong username or password\n\n")
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
				fmt.Scanln()
			}
		}

		tasks, err := api.GetTasks()
		if err != nil {
			log.Fatal(err)
		}
	mainloop:
		for {
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
						fmt.Print("ID has to be of type int\n\n")
						fmt.Scanln()

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
			case "get":
				var id int

				if len(inputSlice) == 2 {
					id, err = strconv.Atoi(inputSlice[1])
					if err != nil {
						fmt.Print("ID has to be of type int\n\n")
						fmt.Scan()

						continue mainloop
					}
				} else {
					fmt.Print("ID: ")
					fmt.Scan(&id)
				}

				currentTask := tasks[id-1]

				frontend.DrawOneTask(id, currentTask)
				fmt.Scanln()
			case "logout":
				break mainloop
			case "exit":
				break fullLoop
			default:
				fmt.Printf("%s not found\n", command)
			}
		}
	}
}
