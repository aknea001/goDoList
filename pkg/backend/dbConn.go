package backend

import (
	"fmt"
	"time"
)

func Connect(done chan bool) {
	fmt.Println("Connecting to db..")
	time.Sleep(3 * time.Second)
	done <- true
}
