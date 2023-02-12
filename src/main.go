package main

import (
	"fmt"

	v1 "github.com/fauzanlucky/consumer-kyc/src/receivers/v1"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	v1.Start()
	fmt.Println("Server stopped")
}
