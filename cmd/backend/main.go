package main

import (
	"fmt"
	"time"

	"github.com/aknea001/goDoList/pkg/backend"
)

func main() {
	done := make(chan bool)

	go backend.Connect(done)

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

outer:
	for {
		select {
		case <-done:
			fmt.Print("\nSuccess!\n")
			break outer
		case <-ticker.C:
			fmt.Print(".")
		}
	}
	fmt.Println("Executing query on db..")
	fmt.Println("NOO it didnt work 🙏")
}
