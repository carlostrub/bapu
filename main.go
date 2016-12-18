package main

import (
	"fmt"
	"os"
)

func main() {

	apiKey := os.Getenv("GANDI_KEY")

	fmt.Println(apiKey)

}
