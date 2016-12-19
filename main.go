package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kolo/xmlrpc"
)

func main() {

	apiKey := os.Getenv("GANDI_KEY")

	//	api, err := xmlrpc.NewClient("https://rpc.gandi.net/xmlrpc/", nil)
	api, err := xmlrpc.NewClient("https://rpc.ote.gandi.net/xmlrpc/", nil)
	if err != nil {
		log.Fatal(err)
	}

	//	var result struct {
	//		Count int `xmlrpc:"count"`
	//	}

	// Count number of instances
	var paasCount *int
	err = api.Call("paas.count", apiKey, &paasCount)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(*paasCount)

}
