package main

import (
	"fmt"

	"github.com/rafamrslima/distributor/internal/db"
	"github.com/rafamrslima/distributor/internal/messaging"
)

func main() {
	fmt.Println("hello project")

	go func() {
		err := messaging.StartMessageListener()

		if err != nil {
			return
		}
	}()

	_, err := db.Connect()
	if err != nil {
		return
	}
}
