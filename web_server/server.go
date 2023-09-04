package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"strconv"
	"strings"
)

type Message struct {
	Type   string
	Body   []byte
	Route  string
	Method string
}

func main() {
	// Listen for incoming connections on a specific port
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		// Accept incoming connection
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Handle the connection (e.g., read and write data)
	// For simplicity, we'll just echo back what we receive

	message, err := readMessage(conn)
	if err != nil {
		fmt.Println("Error in reading message:", err)
		return
	}

	switch {
	case message.Route == "/hello" && message.Method == "GET":
		handleHelloRoute(conn)
	case message.Route == "/bye" && message.Method == "POST":
		handleByeRoute(conn, message.Body)
	default:
		handleUnknownRoute(conn, message.Route, message.Method)
	}

	// var message Message
	// if err := json.Unmarshal([]byte(recivedData), &message); err != nil {
	//     fmt.Println("Error unmarshalling json :", err)
	//     return
	// }
	// fmt.Printf("Received: ", message.Text)
}

func readMessage(conn net.Conn) (Message, error) {
	var message Message
	headerBytes := make([]byte, 2048)
	n, err := conn.Read(headerBytes)
	if err != nil {
		fmt.Println("Error in reading first 10 bytes :", err)
		return Message{}, err
	}
	dataRecived := string(headerBytes[:n])
	// fmt.Println("Received: ", dataRecived)

	protocolAndRouteRaw := strings.Split(dataRecived, "\r\n")
	fmt.Println("Received: ", protocolAndRouteRaw, dataRecived)
	for i := 0; i < len(protocolAndRouteRaw); i++ {
		currentHead := strings.Split(protocolAndRouteRaw[i], ": ")
		if len(currentHead) > 1 {
			fmt.Println("I: ", i, " head: ", currentHead[0], " value: ", currentHead[1])
		} else {
			fmt.Println("I: ", i, " head: ", currentHead[0])
		}
	}
	protocolAndRoute := strings.Split(protocolAndRouteRaw[0], " ")
	// fmt.Println("FDSfsdfF", protocolAndRouteRaw, len(protocolAndRouteRaw), dataRecived[len(dataRecived)-21:])
	messageMethod := protocolAndRoute[0]
	route := protocolAndRoute[1]
	var payload []byte
	if messageMethod != "GET" {
		// paylodRaw := strings.Split(protocolAndRouteRaw[1], "\\n")
		// fmt.Println("FDSF", paylodRaw, len(paylodRaw))
		payload = []byte(dataRecived[21:])
	}

	// Now you have extracted the message type, route, and payload
	message.Method = string(messageMethod)
	message.Type = "request"
	message.Route = route
	message.Body = payload
	return message, nil
}

func handleHelloRoute(conn net.Conn) {
	filePath := "www/hello.html"

	// Read the file contents into a byte slice
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// Convert the byte slice to a string
	fileContentStr := string(fileContent)
	contentType := "text/html"
	sendResponse(conn, fileContentStr, contentType)

}

func handleByeRoute(conn net.Conn, body []byte) {
}

func handleUnknownRoute(conn net.Conn, route string, method string) {
	// errorMessage := "Route not found " + route
	response := "<p>Unknown route: " + route + "</p><p>Method: " + method + "</p>"
	contentType := "text/html"
	// Build the response with headers
	sendResponse(conn, response, contentType)
}

func sendResponse(conn net.Conn, response string, contentType string) {
	contentLength := strconv.Itoa(len(response))
	responseWithHeaders := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: " + contentType + "\r\n" +
		"Content-Length: " + contentLength + "\r\n" +
		"\r\n" +
		response
	// buffer := []byte(errorMessage)
	_, err := conn.Write([]byte(responseWithHeaders))
	if err != nil {
		fmt.Println("Error writing in route not found :", err)
	}
}
