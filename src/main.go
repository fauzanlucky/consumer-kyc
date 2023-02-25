package main

import (
	"fmt"

	v1 "github.com/forkyid/consumer-kyc-update/src/receivers/v1"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	v1.Start()
	fmt.Println("Server stopped")
}
