package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"concurrency_go_course/internal/network"
)

var address string

func init() {
	flag.StringVar(&address, "addr", "127.0.0.1:3223", "database server address")
}

func main() {
	flag.Parse()

	client, err := network.NewClient(address)
	if err != nil {
		fmt.Println("Error starting client:", err.Error())
		os.Exit(1)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter request:")
	for {
		request, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Printf("Error reading request: %v\n", err)
			continue
		}

		resp, err := client.Send(request)
		if err != nil {
			fmt.Printf("unable to send request: %v\n", err)
			continue
		}

		fmt.Println("Server response: ", string(resp))
		fmt.Println("Enter request:")
	}
}
