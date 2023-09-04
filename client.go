package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Message struct {
	Text string `json:"text"`
}

func main() {
	// Connect to the server
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer conn.Close()
	maxAttempts := 2
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		message := Message{
			Text: "Hello, Server!",
		}
		jsonStr, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error marshalling :", err)
			return
		}
		_, err = conn.Write([]byte(jsonStr))
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		fmt.Println("Sent json: ", string(jsonStr))
		// // Receive and print the response
		// buffer := make([]byte, 1024)
		// n, err := conn.Read(buffer)
		// if err != nil {
		// 	fmt.Println("Error:", err)
		// 	return
		// }

		// response := buffer[:n]
		// fmt.Printf("Received: %s", response)
	}
	// Send data to the server

}
